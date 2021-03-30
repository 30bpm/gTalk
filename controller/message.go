package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/groomer/gTalk/config"
)

type Test struct {
	Test string `json:"Test"`
}

type RequestMessageForm struct {
	Customer    string `json:"customer"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Notice      string `json:"notice"`
	PhoneNumber string `json:"phone_number"`
}

type MsgResponse struct {
	UID    uint64 `json:"uid"`
	Date   uint64 `json:"date"`
	Status string `json:"status"`
}

func NoticeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	var msg RequestMessageForm

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(err)
		return
	}

	Message := fmt.Sprintf("안녕하세요 %s님 예약해주셔서 감사합니다 ^_^\n\n%s님께서 예약해주신 미용은\n%s %s에 진행될 예정입니다.\n\n%s에 중요한 내용이 들어있으니 꼭 읽어주세요😎\n그럼 늦지않게 %s %s에 뵙겠습니다\n\n추운 날 감기 조심하세요~❣\n",
		msg.Customer,
		msg.Customer,
		msg.Date,
		msg.Time,
		msg.Notice,
		msg.Date,
		msg.Time)

	URL := "/request/kakao.json"
	data, _ := json.Marshal(map[string]string{
		"service":  strconv.Itoa(config.MESSAGE_NOTICE_SERVICE_NUMBER),
		"message":  Message,
		"mobile":   msg.PhoneNumber,
		"template": "10001",
	})
	req, err := http.NewRequest("POST", config.MESSAGE_API_URL+URL, bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return
	}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{Timeout: timeout}

	req.Header.Set("authToken", config.MESSAGE_API_KEY)
	req.Header.Set("serverName", config.MESSAGE_API_ID)
	req.Header.Set("paymentType", "P")

	msg_response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer msg_response.Body.Close()

	body, err := ioutil.ReadAll(msg_response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var resp MsgResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Println(err)
		return
	}
	if resp.Status != "OK" {
		json.NewEncoder(w).Encode(map[string]string{"msg": "fail"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"msg": "success"})
}
