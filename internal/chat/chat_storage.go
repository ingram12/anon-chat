package chat

import (
	"errors"
	"sort"
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
		Messages:  make(map[int]Message),
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

func (s *Storage) GetMessages(chatID int, limit int) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return nil, ErrChatNotFound
	}

	messages := make([]Message, 0, len(chat.Messages))
	for _, msg := range chat.Messages {
		messages = append(messages, msg)
	}

	// Sort messages by ID (chronological order)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].ID < messages[j].ID
	})

	if limit > 0 && len(messages) > limit {
		return messages[len(messages)-limit:], nil
	}
	return messages, nil
}

func (s *Storage) MarkMessageDelivered(chatID, messageID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return ErrChatNotFound
	}

	msg, exists := chat.Messages[int(messageID)]
	if !exists {
		return ErrMessageNotFound
	}

	msg.IsDelivered = true
	chat.Messages[int(messageID)] = msg
	return nil
}

func (s *Storage) GetUserChats(userID [36]byte) ([]*Chat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userChats []*Chat
	for _, chat := range s.chats {
		if chat.UserID1 == userID || chat.UserID2 == userID {
			userChats = append(userChats, chat)
		}
	}

	return userChats, nil
}
