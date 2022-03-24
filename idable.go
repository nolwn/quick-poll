package main

type idable interface {
	id() string
	setId(id string)
}
