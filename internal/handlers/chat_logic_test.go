package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"anon-chat/internal/pow"

	"github.com/labstack/echo/v4"
)

type testContext struct {
	t      *testing.T
	e      *echo.Echo
	server *Server
}

func setupTest(t *testing.T) *testContext {
	configuration := config.NewConfig()
	configuration.FirstChallengeDifficulty = 1

	server := NewServer(configuration)
	e := echo.New()
	return &testContext{
		t:      t,
		e:      e,
		server: server,
	}
}

func TestChatFlow(t *testing.T) {
	tc := setupTest(t)

	response1 := tc.getFirstChallenge()
	solveResp1 := tc.solveFirstChallenge(response1)
	regResp1 := tc.registerUser(solveResp1)

	response2 := tc.getFirstChallenge()
	solveResp2 := tc.solveFirstChallenge(response2)
	regResp2 := tc.registerUser(solveResp2)

	// Wait for chat
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		tc.waitForChat(regResp1.UserId)
	}()
	go func() {
		defer wg.Done()
		tc.waitForChat(regResp2.UserId)
	}()
	wg.Wait()

	// Update chat
	tc.updateChat(regResp1.UserId)
	tc.updateChat(regResp2.UserId)

	// Send messages
	tc.sendChatMessage(regResp1.UserId, "Hello from user 1")
	tc.sendChatMessage(regResp2.UserId, "Hello from user 2")

	updateResp1 := tc.updateChat(regResp1.UserId)
	if len(updateResp1.Messages) == 0 {
		t.Error("Expected messages in update response for user 1")
	}
	updateResp2 := tc.updateChat(regResp2.UserId)
	if len(updateResp2.Messages) == 0 {
		t.Error("Expected messages in update response for user 2")
	}

	tc.quitChat(regResp1.UserId)
	tc.quitChat(regResp2.UserId)

	wg.Add(2)
	go func() {
		defer wg.Done()
		tc.waitForChat(regResp1.UserId)
	}()
	go func() {
		defer wg.Done()
		tc.waitForChat(regResp2.UserId)
	}()
	wg.Wait()

	tc.sendChatMessage(regResp1.UserId, "Hello from user 1")

	updateResp1 = tc.updateChat(regResp1.UserId)
	if len(updateResp1.Messages) != 0 {
		t.Error("Expected zero messages")
	}

	updateResp2 = tc.updateChat(regResp2.UserId)
	if len(updateResp2.Messages) == 0 && updateResp2.Messages[0].Message != "Hello from user 1" {
		t.Error("Expected messages in update response for user 2")
	}
}

func (tc *testContext) getFirstChallenge() api.GetFirstChallengeResponse {
	rec1 := httptest.NewRecorder()
	ctx1 := tc.e.NewContext(httptest.NewRequest("GET", "/", nil), rec1)

	err := GetFirstChallenge(ctx1, tc.server.cfg, tc.server.rotatingToken)
	if err != nil {
		tc.t.Fatalf("Failed to get first challenge: %v", err)
	}

	var response1 api.GetFirstChallengeResponse
	if err := json.Unmarshal(rec1.Body.Bytes(), &response1); err != nil {
		tc.t.Fatalf("Failed to parse first response: %v", err)
	}

	return response1
}

func (tc *testContext) solveFirstChallenge(response api.GetFirstChallengeResponse) api.SolveFirstChallengeResponse {
	nonce, err := pow.SolveChallenge(response.Challenge, int(response.Difficulty))
	if err != nil {
		tc.t.Fatalf("Failed to solve challenge: %v", err)
	}

	solveReq := api.SolveFirstChallengeRequest{
		Challenge:  response.Challenge,
		Token:      response.Token,
		Nonce:      nonce,
		Difficulty: response.Difficulty,
	}

	solveReqJSON, err := json.Marshal(solveReq)
	if err != nil {
		tc.t.Fatalf("Failed to marshal solve request: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", bytes.NewReader(solveReqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := tc.e.NewContext(req, rec)

	err = SolveFirstChallenge(ctx, tc.server.cfg, tc.server.userStorage, tc.server.rotatingToken)
	if err != nil {
		tc.t.Fatalf("Failed to solve first challenge: %v", err)
	}

	var solveResp api.SolveFirstChallengeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &solveResp); err != nil {
		tc.t.Fatalf("Failed to parse solve response: %v", err)
	}

	return solveResp
}

func (tc *testContext) registerUser(solveResp api.SolveFirstChallengeResponse) api.RegisterUserResponse {
	regNonce, err := pow.SolveChallenge(solveResp.Challenge, int(solveResp.Difficulty))
	if err != nil {
		tc.t.Fatalf("Failed to solve registration challenge: %v", err)
	}

	regReq := api.RegisterUserRequest{
		UserId:     solveResp.UserId,
		Challenge:  solveResp.Challenge,
		Nonce:      regNonce,
		Difficulty: solveResp.Difficulty,
		Nickname:   "TestUser",
		PublicKey:  "TestPublicKey",
		Tags:       []string{"test"},
	}

	regReqJSON, err := json.Marshal(regReq)
	if err != nil {
		tc.t.Fatalf("Failed to marshal registration request: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", bytes.NewReader(regReqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := tc.e.NewContext(req, rec)

	err = RegisterUser(ctx, tc.server.userStorage)
	if err != nil {
		tc.t.Fatalf("Failed to register user: %v", err)
	}

	var regResp api.RegisterUserResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &regResp); err != nil {
		tc.t.Fatalf("Failed to parse registration response: %v", err)
	}

	return regResp
}

func (tc *testContext) waitForChat(userID string) api.WaitForChatResponse {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := tc.e.NewContext(req, rec)
	ctx.SetParamNames("userId")
	ctx.SetParamValues(userID)

	err := WaitForChat(ctx, userID, tc.server.userStorage, tc.server.chatStorage, tc.server.waitingQueue)
	if err != nil {
		tc.t.Fatalf("Failed to wait for chat: %v", err)
	}

	var resp api.WaitForChatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		tc.t.Fatalf("Failed to parse wait chat response: %v", err)
	}

	return resp
}

func (tc *testContext) updateChat(userID string) api.UpdateChatResponse {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := tc.e.NewContext(req, rec)
	ctx.SetParamNames("userId")
	ctx.SetParamValues(userID)

	err := UpdateChat(ctx, userID, tc.server.userStorage, tc.server.chatStorage)
	if err != nil {
		tc.t.Fatalf("Failed to update chat: %v", err)
	}

	if rec.Code == http.StatusBadRequest {
		tc.t.Errorf("Expected 400 Bad Request, got %d", rec.Code)
	}

	var resp api.UpdateChatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		tc.t.Fatalf("Failed to parse update chat response: %v", err)
	}

	return resp
}

func (tc *testContext) sendChatMessage(userID string, message string) api.SendChatMessageResponse {
	sendReq := api.SendChatMessageRequest{
		Message: message,
	}

	sendReqJSON, err := json.Marshal(sendReq)
	if err != nil {
		tc.t.Fatalf("Failed to marshal send message request: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", bytes.NewReader(sendReqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := tc.e.NewContext(req, rec)
	ctx.SetParamNames("userId")
	ctx.SetParamValues(userID)

	err = SendChatMessage(ctx, userID, tc.server.userStorage, tc.server.chatStorage)
	if err != nil {
		tc.t.Fatalf("Failed to send chat message: %v", err)
	}

	var resp api.SendChatMessageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		tc.t.Fatalf("Failed to parse send message response: %v", err)
	}

	return resp
}

func (tc *testContext) quitChat(userID string) api.QuitChatResponse {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := tc.e.NewContext(req, rec)
	ctx.SetParamNames("userId")
	ctx.SetParamValues(userID)

	err := QuitChat(ctx, userID, tc.server.userStorage, tc.server.chatStorage)
	if err != nil {
		tc.t.Fatalf("Failed to quit chat: %v", err)
	}

	var resp api.QuitChatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		tc.t.Fatalf("Failed to parse quit chat response: %v", err)
	}

	return resp
}
