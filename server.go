package main

import (
	"log"
	"net/http"

	"github.com/nolwn/go-router"
)

func main() {
	port := "3000"

	r := router.Router{}

	r.AddRoute(http.MethodGet, "/polls", getAll)
	r.AddRoute(http.MethodPost, "/polls", addPoll)

	r.AddRoute(http.MethodGet, "/polls/:id", getById)

	r.AddRoute(http.MethodPost, "/polls/:id/votes", addVote)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
