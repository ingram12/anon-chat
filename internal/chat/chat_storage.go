package chat

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrChatNotFound    = errors.New("chat not found")
	ErrMessageNotFound = errors.New("message not found")
)

type Storage struct {
	mu        sync.RWMutex
	chats     map[int]*Chat
	lastID    int
	lastMsgID int
}

func NewChatStorage() *Storage {
	return &Storage{
		chats: make(map[int]*Chat),
	}
}

func (s *Storage) CreateChat(userID1, userID2 [36]byte) (*Chat, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	chat := &Chat{
		ID:        s.lastID,
		CreatedAt: time.Now(),
		UserID1:   userID1,
		UserID2:   userID2,
		Messages:  make([]Message, 0, 10),
	}
	s.chats[chat.ID] = chat
	return chat, nil
}

func (s *Storage) GetChat(chatID int) (*Chat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return nil, ErrChatNotFound
	}
	return chat, nil
}

func (s *Storage) AddMessage(chatID int, userID [36]byte, message string) (*Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return nil, ErrChatNotFound
	}

	// Verify user is part of the chat
	if chat.UserID1 != userID && chat.UserID2 != userID {
		return nil, errors.New("user is not part of this chat")
	}

	s.lastMsgID++
	msg := Message{
		ID:          s.lastMsgID,
		CreatedAt:   time.Now(),
		UserID:      userID,
		Message:     message,
		IsDelivered: false,
	}

	chat.Messages[int(s.lastMsgID)] = msg
	return &msg, nil
}
