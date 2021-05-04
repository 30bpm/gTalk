package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/groomer/gTalk/config"
)

type Ping struct {
	Ping string
}

func Pong(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Ping{"Pong"})
}

var MessageStopper map[uint](chan bool)

func IPTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	SalonID := "1"
	fmt.Println("salon ID: ", SalonID)
	postBody, _ := json.Marshal(map[string]string{"salon": SalonID})
	resp, err := http.Post(config.GROOMER_NOTE_URL+"/alarmtalk/decrease/", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	sb := string(body)
	fmt.Println(sb)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "success"})
}
