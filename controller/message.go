package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/groomer/gTalk/config"
	"github.com/groomer/gTalk/db"
)

type Test struct {
	Test string `json:"Test"`
}

type RequestMessageForm struct {
	Phone    string `json:"phone"`
	Customer string `json:"customer"`
	Template string `json:"template"`
	Notice   string `json:"notice"`
	Salon    uint64 `json:"salon"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	Booking  uint64 `json:"booking"`
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

	if err := MessageCountDecrease(fmt.Sprint(msg.Salon)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"msg": "credit is zero"})
		return
	}

	Text := fmt.Sprintf("안녕하세요 %s님 예약해주셔서 감사합니다 ^_^\n\n%s님께서 예약해주신 미용은\n%s %s에 진행될 예정입니다.\n\n%s에 중요한 내용이 들어있으니 꼭 읽어주세요😎\n그럼 늦지않게 %s %s에 뵙겠습니다\n\n추운 날 감기 조심하세요~❣\n",
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
		"message":  Text,
		"mobile":   msg.Phone,
		"template": msg.Template,
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

	var message db.Message
	message.Phone = msg.Phone
	message.Template = msg.Template
	message.Salon = msg.Salon
	message.Date = msg.Date
	message.Time = msg.Time
	message.Done = true
	message.Text = Text
	db.DB.Create(&message)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "success"})
}

func Reminder(message db.Message) bool {
	HM := strings.Split(message.Time, ":")
	now := time.Now()
	hour, _ := strconv.Atoi(HM[0])
	hour = hour - now.Hour()
	minute, _ := strconv.Atoi(HM[1])
	minute = minute - now.Minute()

	date, _ := time.Parse("2006-01-02", message.Date)
	nowDate, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	if !(date.Year() == nowDate.Year() && date.Month() == nowDate.Month() && date.Day() == nowDate.Day()) {
		return true
	}
	if minute < 0 && hour > 0 {
		hour -= 1
		minute += 60
	}

	if minute < 0 || hour < 0 {
		return false
	}

	var duration time.Duration
	if hour == 0 && minute == 0 {
		duration = 1 * time.Minute
	} else {
		duration = time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute
	}

	ticker := time.NewTicker(duration)
	MessageStopper[message.Booking] = make(chan bool)
	fmt.Println("map: ", MessageStopper)
	fmt.Println(message.ID, time.Now().Add(duration))
	for {
		select {
		case <-ticker.C:
			fmt.Println("send")
			URL := "/request/kakao.json"
			data, _ := json.Marshal(map[string]string{
				"service":  strconv.Itoa(config.MESSAGE_NOTICE_SERVICE_NUMBER),
				"message":  message.Text,
				"mobile":   message.Phone,
				"template": message.Template,
			})
			if err := MessageCountDecrease(fmt.Sprint(message.Salon)); err != nil {
				return false
			}
			req, err := http.NewRequest("POST", config.MESSAGE_API_URL+URL, bytes.NewBuffer(data))
			if err != nil {
				log.Println(err)
				return false
			}
			timeout := time.Duration(5 * time.Second)
			client := http.Client{Timeout: timeout}

			req.Header.Set("authToken", config.MESSAGE_API_KEY)
			req.Header.Set("serverName", config.MESSAGE_API_ID)
			req.Header.Set("paymentType", "P")

			msg_response, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return false
			}

			defer msg_response.Body.Close()

			body, err := ioutil.ReadAll(msg_response.Body)
			if err != nil {
				log.Println(err)
				return false
			}

			var resp MsgResponse
			if err := json.Unmarshal(body, &resp); err != nil {
				log.Println(err)
				return false
			}
			if resp.Status != "OK" {
				return false
			}
			if message.Done == false {
				message.Done = true
				db.DB.Save(&message)
			}
			delete(MessageStopper, message.Booking)
			return true
		case <-MessageStopper[message.Booking]:
			fmt.Println("map: ", MessageStopper)
			fmt.Println("[Log] ", message.Booking, "stopped")
			close(MessageStopper[message.Booking])
			delete(MessageStopper, message.Booking)
			db.DB.Unscoped().Delete(&message)
			fmt.Println("map: ", MessageStopper)
			return false
		}
	}
}

func GetReminders() bool {
	fmt.Println("get reminders time: ", time.Now())
	var messages []db.Message
	date := time.Now().Format("2006-01-02")
	db.DB.Where("date = ?", date).Find(&messages)
	fmt.Println("time: ", time.Now())
	fmt.Println("len: ", len(messages))
	for _, message := range messages {
		fmt.Println(message.ID, message.Date, message.Time)
		go Reminder(message)
	}
	return true
}

func PostReminder(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(writer, "invalid_http_method")
		return
	}
	var msg RequestMessageForm

	if err := json.NewDecoder(request.Body).Decode(&msg); err != nil {
		log.Println(err)
		return
	}

	Text := fmt.Sprintf("안녕하세요 %s님 예약해주셔서 감사합니다 ^_^\n\n%s님께서 예약해주신 미용은\n%s %s에 진행될 예정입니다.\n\n%s에 중요한 내용이 들어있으니 꼭 읽어주세요😎\n그럼 늦지않게 %s %s에 뵙겠습니다\n\n추운 날 감기 조심하세요~❣\n",
		msg.Customer,
		msg.Customer,
		msg.Date,
		msg.Time,
		msg.Notice,
		msg.Date,
		msg.Time)

	var message db.Message
	message.Phone = msg.Phone
	message.Template = msg.Template
	message.Salon = msg.Salon
	message.Date = msg.Date
	message.Time = msg.Time
	message.Text = Text
	message.Booking = msg.Booking
	message.Done = false
	if err := db.DB.Create(&message).Error; err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode((map[string]string{"msg": fmt.Sprint(err)}))
		return
	}

	go Reminder(message)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"msg": "success", "id": fmt.Sprint(message.ID)})
}

type MessageBooking struct {
	Booking uint64 `jons:"booking"`
}

func DeleteReminder(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(writer, "invalid_http_method")
		return
	}
	var messageBooking MessageBooking

	if err := json.NewDecoder(request.Body).Decode(&messageBooking); err != nil {
		log.Println(err)
		return
	}
	if channel, ok := MessageStopper[messageBooking.Booking]; ok {
		channel <- true
	}
	// delete(MessageStopper, messageBooking.Booking)

	json.NewEncoder(writer).Encode(map[string]string{"msg": "success"})
}
