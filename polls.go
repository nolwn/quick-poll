package main

// newPoll contains the fields a user can send to create a new poll
type newPoll struct {

	// title of the poll
	Title string

	// values of the options that can be voted on
	Options []string
}

// poll contains all the data for a poll
type poll struct {

	// string representation of the poll's id
	Id string `bson:",omitempty"`

	// title of the poll
	Title string

	// list of pollOptions which can be voted on
	Options []pollOption
}

// pollOptions is a list of values that a user can vote on. It also has Votes which is a
// cache that stores that number of users who have voted for that option.
type pollOption struct {

	// option value that the user will see and can vote on
	Value string

	// count of all the users that have voted on for this item
	Votes int
}

// id returns the item's id. It helps fullfil the idable interface.
func (p *poll) id() string {
	return p.Id
}

// setId takes a string and sets the id based on that string. It helps fullfil the idable
// interface.
func (p *poll) setId(id string) {
	p.Id = id
}
