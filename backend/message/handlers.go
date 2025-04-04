package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader to convert HTTP to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { //CheckOrigin: func(r *http.Request) bool { return true } allows connections from any origin.
		return true
	},
}

// Client represents a WebSocket connection
type Client struct {
	ID   string
	Conn *websocket.Conn //websocket connection b/w user and server
	Send chan []byte     //channel for outgoing messages
}

// MessageHub --> to manage all active websocket connection
/*
Clients: Maps user IDs to WebSocket connections.
Groups: Maps group IDs to sets of clients
*/
type MessageHub struct {
	Clients       map[string]*Client          // Tracks online users (userID -> Client)
	Groups        map[string]map[*Client]bool // Tracks group members (groupID -> Clients)
	Broadcast     chan MessagePayload
	Register      chan *Client
	Unregister    chan *Client
	DirectMessage chan MessagePayload
}

// payload --> to hold message information
type MessagePayload struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id,omitempty"`
	GroupID    string `json:"group_id,omitempty"`
	Content    string `json:"content"`
}

var Hub = MessageHub{
	Clients:       make(map[string]*Client),
	Groups:        map[string]map[*Client]bool{},
	Broadcast:     make(chan MessagePayload),
	Register:      make(chan *Client),
	Unregister:    make(chan *Client),
	DirectMessage: make(chan MessagePayload),
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil) // returns a websocket connection and error --> func (u *websocket.Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)

	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// Getting user ID from query parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("User ID missing")
		conn.Close()
		return
	}

	//create a new client
	client := &Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan []byte),
	}

	//register that client into Hub

	Hub.Register <- client

	// Start goroutines for reading and writing messages
	go client.readMessages()
	go client.writeMessages()

}

/*RemoteAddr returns the WebSocket location for the connection for client, or the Websocket Origin for server.
func (ws *Conn) RemoteAddr() net.Addr {
	if ws.IsClientConn() {
		return &Addr{ws.config.Location}
	}
	return &Addr{ws.config.Origin}
}
*/

func StartHub() {
	for {
		select {
		case client := <-Hub.Register: //receive from Register channel in Hub
			//register the user ---> when they connect
			userID := client.Conn.RemoteAddr().String()
			Hub.Clients[userID] = client
			log.Printf("User %s connected", userID)

		case client := <-Hub.Unregister:
			userID := client.Conn.RemoteAddr().String()
			_, exists := Hub.Clients[userID]
			if exists {
				delete(Hub.Clients, userID) // delete the client.conn address
				close(client.Send)          //close the send channel corresponding to that channel
				log.Printf("User %s disconnected", userID)
			}

		// P2P messaging
		case msg := <-Hub.DirectMessage: //DirectMessage is MessagePayload channel, same with broadcast
			recepient, exists := Hub.Clients[msg.ReceiverID]
			if exists {
				messageJSON, _ := json.Marshal(msg)
				recepient.Send <- messageJSON //send msg to recepient's send channels
			}

		//Group chat
		case msg := <-Hub.Broadcast:
			recepients, exist := Hub.Groups[msg.GroupID]
			if exist {
				messageJSON, _ := json.Marshal(msg)
				for recepient := range recepients {
					recepient.Send <- messageJSON
				}
			}
		}
	}

}

// client writes message
func (c *Client) writeMessages() {
	defer c.Conn.Close()

	for msg := range c.Send { // Only read messages meant for this client
		err := c.Conn.WriteMessage(websocket.TextMessage, msg) //func (c *websocket.Conn) WriteMessage(messageType int, data []byte) error
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}

}

func (c *Client) readMessages() {
	defer c.Conn.Close()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var payload MessagePayload
		// Parse message JSON
		err = json.Unmarshal(message, &payload)

		if err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		if payload.GroupID != "" {
			Hub.Broadcast <- payload
		} else {
			Hub.DirectMessage <- payload
		}
	}
}

/*

Peer to peer
-----------------------------------------------
Alice sends the message over WebSocket.
The server receives it and puts it in Hub.DirectMessage.
StartHub picks it up and looks for Bob in Hub.Clients.
The message is sent only to Bob’s Send channel`.
Bob’s writeMessages function reads the message and sends it over WebSocket.


Group messaging
-------------------------------------------------------------------
Alice’s message is added to Hub.Broadcast.
StartHub checks Hub.Groups["G1"].
Bob and Charlie receive the message in their Send channels.
Their writeMessages functions send the message over WebSocket.
*/

func AddUserToGroup(userID, groupID string) {

	client, exists := Hub.Clients[userID]

	if !exists {
		log.Printf("User %s not found", userID)
		return
	}

	//checking if the group exists
	_, exist := Hub.Groups[groupID]

	if !exist {
		Hub.Groups[groupID] = make(map[*Client]bool)
	}

	Hub.Groups[groupID][client] = true
	log.Printf("User %s added to group %s", userID, groupID)

}

func RemoveUserFromGroup(userID, groupID string) {

	client, exists := Hub.Clients[userID]
	if !exists {
		log.Printf("User %s not found", userID)
		return
	}

	group, groupExists := Hub.Groups[groupID]
	if !groupExists {
		log.Printf("Group %s not found", groupID)
		return
	}

	exist := group[client]
	if exist {
		delete(group, client)
		log.Printf("User %s removed from Group %s", userID, groupID)
	} else {
		log.Printf("user %s is not present in the group %s ", userID, groupID)
	}
}
