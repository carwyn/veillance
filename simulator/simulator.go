package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/antage/eventsource"
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

func main() {
	es := eventsource.New(
		&eventsource.Settings{
			Timeout:        5 * time.Second,
			CloseOnTimeout: false,
			IdleTimeout:    30 * time.Minute,
		},
		func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
				[]byte("Content-Type: text/event-stream"),
				[]byte("Cache-Control: no-cache"),
				[]byte("Connection: keep-alive"),
			}
		},
	)

	defer es.Close()
	http.Handle("/events", es)
	go func() {
		i, id := 0, 0
		bl := len(bricksong)
		dl := len(domains)
		for {
			bi := i % bl
			di := i % dl
			f := Fragment{bi, Source{Name: "SimUser", Type: "User"}, bricksong[bi]}
			d := Domain{di, Source{Name: "SimUser", Type: "User"}, domains[di]}
			jsb, _ := json.Marshal(f)
			jsd, _ := json.Marshal(d)
			es.SendEventMessage(string(jsb), "fragment", strconv.Itoa(id))
			id++
			es.SendEventMessage(string(jsd), "domain", strconv.Itoa(id))
			id++

			if i > bl {
				s := Selection{rand.Intn(bl), []int{1, 3}}
				jss, _ := json.Marshal(s)
				es.SendEventMessage(string(jss), "selection", strconv.Itoa(id))
				id++
			}
			i++
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
