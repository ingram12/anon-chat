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
	Mu     sync.RWMutex
	chats  map[int]*Chat
	lastID int
}

func NewChatStorage() *Storage {
	return &Storage{
		chats: make(map[int]*Chat),
	}
}

func (s *Storage) CreateChat(userID1, userID2 string) (*Chat, error) {
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
	chat, exists := s.chats[chatID]
	if !exists {
		return nil, ErrChatNotFound
	}
	return chat, nil
}

func (s *Storage) GetPeerMessages(chatID int, userID string) ([]Message, error) {
	chat, exists := s.chats[chatID]
	if !exists {
		return []Message{}, ErrChatNotFound
	}

	return chat.GetPeerMessages(userID), nil
}

func (s *Storage) HasNewMessages(chatID int, userID string) bool {
	chat, exists := s.chats[chatID]
	if !exists {
		return false
	}
	if chat.UserID1 == userID {
		return len(chat.User2Messages) > 0
	} else if chat.UserID2 == userID {
		return len(chat.User1Messages) > 0
	}
	return false
}

func (s *Storage) RemovePeerMessages(chatID int, userID string) error {
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

func (s *Storage) AddMessage(chatID int, userID string, message string) (time.Time, error) {
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

func (s *Storage) QuitChat(chatID int, userID string) error {
	chat, exists := s.chats[chatID]
	if !exists {
		return ErrChatNotFound
	}

	if chat.UserID1 == userID {
		chat.UserID1 = ""
	} else if chat.UserID2 == userID {
		chat.UserID2 = ""
	} else {
		return errors.New("user not in chat")
	}

	return nil
}

func (s *Storage) IsUserInChat(chatID int, userID string) bool {
	chat, exists := s.chats[chatID]
	if !exists {
		return false
	}

	return chat.IsUserInChat(userID)
}

func (s *Storage) IsActiveChat(chatID int) bool {
	chat, exists := s.chats[chatID]
	if !exists {
		return false
	}

	return chat.IsActive()
}
