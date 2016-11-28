package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxId int

// Chat server.
type Server struct {
	pattern   string
	messages  []Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan Message
	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			//s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

			// broadcast message for all clients
		case msg := <-s.sendAllCh:
			js, _ := json.Marshal(msg)
			log.Println("Send all:", string(js))
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func simulate(server *Server) {
	i := 0
	bl := len(bricksong)
	dl := len(domains)
	for {
		bi := i % bl
		di := i % dl
		f := &Fragment{bi, Source{Name: "User1", Type: "User"}, bricksong[bi]}
		d := &Domain{di, Source{Name: "User2", Type: "User"}, domains[di]}
		//jsb, _ := json.Marshal(f)
		//jsd, _ := json.Marshal(d)
		server.SendAll(f)
		server.SendAll(d)
		//server.SendAll(&Message{"Bob", string(jsb)})
		//server.SendAll(&Message{"Bob", string(jsd)})

		if i > bl {
			s := &Selection{rand.Intn(bl), []int{1, 3}}
			//jss, _ := json.Marshal(s)
			//server.SendAll(&Message{"Bob", string(jss)})
			server.SendAll(s)
		}
		i++
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func startServer() *Server {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := NewServer("/entry")
	go server.Listen()

	//go simulate(server)
	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	go http.ListenAndServe(":8080", nil)

	return server
}
