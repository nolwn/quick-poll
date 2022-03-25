package main

// newVote is all the information a user can pass to vote on a poll
type newVote struct {
	Choice int `json:"choice"`
}

// vote is all the information about a users vote
type vote struct {

	// string representation of the vote's id
	Id string `json:"id"`

	// choice is the number of the choice selected
	Choice int `json:"choice"`

	// poll id is the id of the poll being voted on
	PollId string `json:"pollId"`
}

// id returns the item's id. It helps fullfil the idable interface.
func (v *vote) id() string {
	return v.Id
}

// setId takes a string and sets the id based on that string. It helps fullfil the idable
// interface.
func (v *vote) setId(id string) {
	v.Id = id
}
