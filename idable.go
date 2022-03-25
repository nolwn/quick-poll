package main

// idable should be used for anything that has an id
// in the database.
type idable interface {
	// id should return the id
	id() string

	// setId should take a new id and set it as the item's id
	setId(id string)
}
