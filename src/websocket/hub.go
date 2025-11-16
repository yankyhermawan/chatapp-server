package websocket

import (
	"chatapp/src/dto"
	"chatapp/src/message"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Db         *gorm.DB
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Db:         db,
	}
}

func emitSpecificUser(client *Client, mess dto.ReceivedMessage, h *Hub) {
	if (client.UserId == mess.ReceiverId) || (client.UserId == mess.SenderId) {
		messageToBeSend := dto.SendMessage{
			Message:    mess.Message,
			SenderId:   mess.SenderId,
			ReceiverId: mess.ReceiverId,
		}
		if (client.WsType == "chat") && (client.UserId == mess.SenderId) {
			message.CreateMessage(h.Db, mess)
		}
		jsonBytes, _ := json.Marshal(messageToBeSend)
		fmt.Printf("%s\n", jsonBytes)
		select {
		case client.Send <- jsonBytes:
		default:
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			var parsedMesage dto.ReceivedMessage
			json.Unmarshal(message, &parsedMesage)
			for client := range h.Clients {
				emitSpecificUser(client, parsedMesage, h)
			}
		}
	}
}
