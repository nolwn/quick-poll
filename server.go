package main

import (
	"log"
	"net/http"

	"github.com/nolwn/go-router"
	"github.com/nolwn/quick-poll/polls"
)

func main() {
	port := "3000"

	r := router.Router{}

	r.AddRoute(http.MethodGet, "/polls", polls.GetAll)
	r.AddRoute(http.MethodPost, "/polls", polls.AddPoll)
	r.AddRoute(http.MethodGet, "/polls/:id", polls.GetById)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
