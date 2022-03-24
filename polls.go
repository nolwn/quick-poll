package main

type newPoll struct {
	Title   string
	Options []string
}

type poll struct {
	Id      string `bson:",omitempty"`
	Title   string
	Options []pollOption
}

type pollOption struct {
	Value string
	Votes int
}

func (p *poll) id() string {
	return p.Id
}

func (p *poll) setId(id string) {
	p.Id = id
}
