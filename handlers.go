package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nolwn/go-router"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	optionsMin = 2
	optionsMax = 16
)

type createResponse struct {
	Id string `json:"id"`
}

func getAll(w http.ResponseWriter, r *http.Request) {
	polls, err := query(TablePoll, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(polls)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func getById(w http.ResponseWriter, r *http.Request) {
	params := router.PathParams(r)

	poll, err := queryById(TablePoll, params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if poll == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(poll)

	if err != nil {
		fmt.Printf("Could not marshall poll: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func addPoll(w http.ResponseWriter, r *http.Request) {
	var newPoll newPoll
	var poll poll
	var badRequest bool
	var res []byte

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newPoll)

	if err != nil {
		fmt.Println(err)
		badRequest = true
	} else if len(newPoll.Options) < optionsMin {
		fmt.Printf("Recieved a poll with fewer than %d options", optionsMin)
		badRequest = true
	} else if len(newPoll.Options) > optionsMax {
		fmt.Printf("Recieved a poll with more than %d options", optionsMax)
		badRequest = true
	}

	if badRequest {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Adding a poll requires title and options."))
		return
	}

	poll.Title = newPoll.Title
	poll.Options = make([]pollOption, len(newPoll.Options))

	for i, o := range newPoll.Options {
		poll.Options[i] = pollOption{
			Value: o,
		}
	}

	id, err := add(TablePoll, poll)
	if err == nil {
		res, err = json.Marshal(createResponse{id})
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func addVote(w http.ResponseWriter, r *http.Request) {
	var addVote newVote
	var vote vote
	var res []byte
	var poll poll

	params := router.PathParams(r)
	pollId := params["id"]
	bsonPoll, err := queryById(TablePoll, pollId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mPoll, err := bson.Marshal(bsonPoll)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bson.Unmarshal(mPoll, &poll)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addVote.Choice = -1

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&addVote)

	if err != nil || addVote.Choice < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Provide a choice to make a vote"))
		return
	}

	vote.Choice = addVote.Choice
	vote.PollId = pollId
	poll.Options[addVote.Choice].Votes++

	id, err := add(TableVote, vote)
	if err == nil {
		res, err = json.Marshal(createResponse{id})
	}

	updateErr := update(TablePoll, poll.Id, &poll)
	if err != nil || updateErr != nil {
		fmt.Printf("Didn't Update: %s %s", err, updateErr)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
