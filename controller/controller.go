package controller

import (
	"encoding/json"
	"net/http"
)

type Ping struct {
	Ping string
}

func Pong(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Ping{"Pong"})
}

var MessageStopper map[uint](chan bool)
