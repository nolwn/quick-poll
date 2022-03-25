package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nolwn/go-router"
	"go.mongodb.org/mongo-driver/bson"
)

const (

	// minimum number of options that a poll can have
	optionsMin = 2

	// maximum number of options that a poll can have
	optionsMax = 16
)

// struct for a create item response
type createResponse struct {
	Id string `json:"id"`
}

// struct for an error response
type errorResponse struct {
	Error string `json:"error"`
}

// getAllPolls handles returning all the polls to the caller
func getAllPolls(w http.ResponseWriter, r *http.Request) {
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

// getPollById handles returning the poll whose id corresponds to the id sent
// by the caller.
func getPollById(w http.ResponseWriter, r *http.Request) {
	params := router.PathParams(r)

	poll, err := queryById(TablePoll, params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if poll is nil, we must not have found a poll
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

// addPoll handles creating a new poll based on the infomation passed by the caller
func addPoll(w http.ResponseWriter, r *http.Request) {
	var newPoll newPoll
	var poll poll
	var badRequest bool
	var res []byte
	var errMsg string

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newPoll)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	} else if len(newPoll.Options) < optionsMin {
		errMsg = fmt.Sprintf("Recieved a poll with fewer than %d options", optionsMin)
		res, err = json.Marshal(errorResponse{errMsg})
		badRequest = true

	} else if len(newPoll.Options) > optionsMax {
		errMsg = fmt.Sprintf("Recieved a poll with more than %d options", optionsMax)
		res, err = json.Marshal(errorResponse{errMsg})
		badRequest = true
	}

	// unexpected error exists
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if badRequest {
		w.WriteHeader(http.StatusBadRequest)

		if res != nil {
			w.Write(res)
		} else {
			// something was wrong with the resquest
			w.Write([]byte("Adding a poll requires title and options."))
		}

		return
	}

	poll.Title = newPoll.Title
	poll.Options = make([]pollOption, len(newPoll.Options))

	// fill the poll options array with provided poll values
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

// addVote handles adding the callers vote to the database and updating the vote count
// on the poll
func addVote(w http.ResponseWriter, r *http.Request) {
	var addVote newVote
	var vote vote
	var res []byte
	var poll poll

	params := router.PathParams(r)
	pollId := params["id"]

	// find the poll the caller is voting on
	bsonPoll, err := queryById(TablePoll, pollId)

	if bsonPoll != nil { // no poll with that id was found
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil { // some unexpected error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// to turn a BSON object into a struct, you have to marshal and unmarshal it.
	//
	// marshal into byte array
	mPoll, err := bson.Marshal(bsonPoll)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// unmarshall into the poll object
	err = bson.Unmarshal(mPoll, &poll)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// new structs have zeroed out properties so, to make sure the user has explicity
	// entered a choice, since zero is a valie choice, Choice must be set to something
	// that couldn't be valid.
	addVote.Choice = -1

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&addVote)
	if err != nil || addVote.Choice < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Provide a choice to make a vote"))
		return
	}

	// load vote struct
	vote.Choice = addVote.Choice
	vote.PollId = pollId

	// update vote count on the poll
	poll.Options[addVote.Choice].Votes++

	// create vote in db
	id, err := add(TableVote, vote)
	if err == nil {
		res, err = json.Marshal(createResponse{id})
	}

	// update poll with new vote count
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
