package internal

import (
	"sync"

	"github.com/temuka-api-service/models"
	"gorm.io/gorm"
)

type Hub struct {
	Conversations map[int]*models.Conversation
	Register      chan *models.Participant
	Unregister    chan *models.Participant
	Broadcast     chan *models.Message
	mu            sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Conversations: make(map[int]*models.Conversation),
		Register:      make(chan *models.Participant),
		Unregister:    make(chan *models.Participant),
		Broadcast:     make(chan *models.Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case p := <-h.Register:
			h.mu.Lock()
			conversation, ok := h.Conversations[p.ConversationID]
			if !ok {
				conversation = &models.Conversation{
					ID:           p.ConversationID,
					Participants: []models.Participant{},
				}
				h.Conversations[p.ConversationID] = conversation
			}
			conversation.Participants = append(conversation.Participants, *p)
			h.mu.Unlock()
		case p := <-h.Unregister:
			h.mu.Lock()
			if conversation, ok := h.Conversations[p.ConversationID]; ok {
				for i, participant := range conversation.Participants {
					if participant.ID == p.ID {
						conversation.Participants = append(conversation.Participants[:i], conversation.Participants[i+1:]...)
						break
					}
				}
				if len(conversation.Participants) == 0 {
					delete(h.Conversations, p.ConversationID)
				}
			}
			h.mu.Unlock()
		case m := <-h.Broadcast:
			h.mu.Lock()
			if conversation, ok := h.Conversations[m.ParticipantID]; ok {
				for _, participant := range conversation.Participants {
					participant.Messages = append(participant.Messages, *m)
				}
			}
			h.mu.Unlock()
		}

	}
}

var db *gorm.DB
var hub *Hub

func Init(h *Hub, database *gorm.DB) {
	hub = h
	db = database
}
