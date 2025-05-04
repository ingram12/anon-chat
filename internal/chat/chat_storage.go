package chat

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrChatNotFound = errors.New("chat not found")
)

type Storage struct {
	mu     sync.RWMutex
	chats  map[int]*Chat
	lastID int
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
		ID:            s.lastID,
		CreatedAt:     time.Now(),
		UserID1:       userID1,
		UserID2:       userID2,
		User1Messages: []Message{},
		User2Messages: []Message{},
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

func (s *Storage) GetPeerMessages(chatID int, userID [36]byte) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return []Message{}, ErrChatNotFound
	}

	messages := chat.GetPeerMessages(userID)

	return messages, nil
}

func (s *Storage) RemovePeerMessages(chatID int, userID [36]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return ErrChatNotFound
	}

	if chat.UserID1 == userID {
		chat.User2Messages = []Message{}
	} else if chat.UserID2 == userID {
		chat.User1Messages = []Message{}
	} else {
		return errors.New("user not in chat")
	}

	return nil
}

func (s *Storage) AddMessage(chatID int, userID [36]byte, message string) (time.Time, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return time.Time{}, ErrChatNotFound
	}

	timeNow := time.Now()
	if chat.UserID1 == userID {
		chat.User1Messages = append(chat.User1Messages, Message{
			Timestamp: timeNow,
			Message:   message,
		})
	} else if chat.UserID2 == userID {
		chat.User2Messages = append(chat.User2Messages, Message{
			Timestamp: timeNow,
			Message:   message,
		})
	} else {
		return time.Time{}, errors.New("user not in chat")
	}

	return timeNow, nil
}

func (s *Storage) QuitChat(chatID int, userID [36]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return ErrChatNotFound
	}

	if chat.UserID1 == userID {
		chat.UserID1 = [36]byte{}
	} else if chat.UserID2 == userID {
		chat.UserID2 = [36]byte{}
	} else {
		return errors.New("user not in chat")
	}

	return nil
}

func (s *Storage) IsUserInChat(chatID int, userID [36]byte) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chat, exists := s.chats[chatID]
	if !exists {
		return false
	}

	return chat.IsUserInChat(userID)
}
