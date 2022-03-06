package polls

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nolwn/quick-poll/data"
	"github.com/nolwn/quick-poll/resources"
)

type createResponse struct {
	Id string `json:"id"`
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	polls, err := data.Query(data.TablePoll, resources.AddPoll{})
	if err != nil {
		w.WriteHeader(500)
	}

	b, err := json.Marshal(polls)
	if err != nil {
		w.WriteHeader(500)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func AddPoll(w http.ResponseWriter, r *http.Request) {
	var poll resources.AddPoll
	var badRequest bool
	var res []byte

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&poll)

	if err != nil {
		fmt.Println(err)
		badRequest = true
	} else if len(poll.Options) < 2 {
		fmt.Println("Recieved a poll with one or fewer options")
		badRequest = true
	}

	if badRequest {
		w.WriteHeader(400)
		w.Write([]byte("Adding a poll requires title and options."))
		return
	}

	id, err := data.Add(data.TablePoll, poll)
	if err == nil {
		res, err = json.Marshal(createResponse{id})
	}

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(res)
}

// func GetById(w http.ResponseWriter, r *http.Request) {

// }
