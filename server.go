package main

import (
	"log"
	"net/http"

	"github.com/nolwn/go-router"
)

func main() {

	// TODO: this should be pulled of the environment
	port := "3000"

	r := router.Router{}

	// poll endpoints
	r.AddRoute(http.MethodGet, "/polls", getAllPolls)
	r.AddRoute(http.MethodPost, "/polls", addPoll)

	// poll by id endpoints
	r.AddRoute(http.MethodGet, "/polls/:id", getPollById)

	// vote endpoints
	r.AddRoute(http.MethodPost, "/polls/:id/votes", addVote)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
