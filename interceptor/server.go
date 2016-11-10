package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var domains = [...]string{
	"facebook.com",
	"disney.com",
	"google.com",
	"microsoft.com",
	"foobar.io",
	"fast.com",
	"netflix.com",
	"anandtech.com",
	"arstechnica.com",
	"www.bbc.co.uk",
	"tumblr.com",
	"instagram.com",
	"twitter.com",
	"phoronix.com",
	"channel4.com",
	"theregister.co.uk",
	"gov.uk",
}

var bricksong = [...]string{
	"Dear Sir I write this note to inform you of my plight",
	"And at the time of writing I am not a pretty sight",
	"My body is all black and blue, my face a deathly gray",
	"I write this note to tell why Paddy's not at work today",

	"While working on the fourteenth floor, some bricks I had to clear",
	"And to throw them down from off the top seemed quite a good idea",
	"But the gaffer wasn't very pleased, he was an awful sod",
	"He said I had to cart them down the ladder in me hod.",

	"Well clearing all those bricks by hand, it seemed so very slow",
	"So I hoisted up a barrel and secured the rope below",
	"But in my haste to do the job, I was too blind to see",
	"That a barrel full of building bricks is heavier than me.",

	"So when I had untied the rope, the barrel fell like lead",
	"And clinging tightly to the rope I started up instead",
	"I took off like a rocket and to my dismay I found",
	"That half way up I met the bloody barrel coming down.",

	"Well the barrel broke my shoulder as on to the ground it sped",
	"And when I reached the top I banged the pulley with me head",
	"I held on tight, though numb with shock from this almighty blow",
	"And the barrel spilled out half its load fourteen floors below",

	"Now when those building bricks fell from the barrel to the floor",
	"I then outweighed the barrel so I started down once more",
	"I held on tightly to the rope as I flew to the ground",
	"And I landed on those building bricks that were scattered all around.",

	"Now as I lay there on the deck I thought I'd passed the worst",
	"But when the barrel reached the top, that's when the bottom burst",
	"A shower of bricks came down on me, I knew I had no hope",
	"In all of this confusion, I let go the bloody rope.",

	"The barrel being heavier, it started down once more",
	"And landed right on top of me as I lay on the floor",
	"It broke three ribs and my left arm, and I can only say",
	"That I hope you'll understand why Paddy's not at work today.",
}

const channelBufSize = 100

var maxId int = 0

// Chat client.
type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *Message
	doneCh chan bool
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, ws, server, ch, doneCh}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.SendAll(&msg)
			}
		}
	}
}

type Source struct {
	Name string
	Type string
}

type Fragment struct {
	Id int
	Source
	Text string
}

type Selection struct {
	Fid   int
	Words []int
}

type Domain struct {
	Id int
	Source
	Text string
}

type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (self *Message) String() string {
	return self.Author + " says " + self.Body
}

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
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

func (s *Server) SendAll(msg *Message) {
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

func (s *Server) sendAll(msg *Message) {
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
			log.Println("Send all:", msg)
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
		f := Fragment{bi, Source{Name: "SimUser", Type: "User"}, bricksong[bi]}
		d := Domain{di, Source{Name: "SimUser", Type: "User"}, domains[di]}
		jsb, _ := json.Marshal(f)
		jsd, _ := json.Marshal(d)
		server.SendAll(&Message{"Bob", string(jsb)})
		server.SendAll(&Message{"Bob", string(jsd)})

		if i > bl {
			s := Selection{rand.Intn(bl), []int{1, 3}}
			jss, _ := json.Marshal(s)
			server.SendAll(&Message{"Bob", string(jss)})
		}
		i++
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := NewServer("/entry")
	go server.Listen()

	go simulate(server)
	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
