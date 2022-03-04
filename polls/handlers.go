package polls

import (
	"encoding/json"
	"net/http"

	"github.com/nolwn/quick-poll/data"
	"github.com/nolwn/quick-poll/resources"
)

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

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&poll)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Adding a poll requires title and options."))
		return
	}

	id, err := data.Add(data.TablePoll, poll)

	if err != nil {
		w.WriteHeader(500)
	}

	w.WriteHeader(201)
	w.Write([]byte(id))
}

// func GetById(w http.ResponseWriter, r *http.Request) {

// }
