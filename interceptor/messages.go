package main

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

//type Message struct {
//	Author string `json:"author"`
//	Body   string `json:"body"`
//}

type Message interface {
}

//func (self *Message) String() string {
//	return self.Author + " says " + self.Body
//}
