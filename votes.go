package main

type newVote struct {
	Choice int `json:"choice"`
}

type vote struct {
	Id     string `json:"id"`
	Choice int    `json:"choice"`
	PollId string `json:"pollId"`
}

func (v *vote) id() string {
	return v.Id
}

func (v *vote) setId(id string) {
	v.Id = id
}
