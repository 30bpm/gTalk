package router

import (
	"fmt"
	"net/http"

	"github.com/groomer/gTalk/controller"
	"github.com/groomer/gTalk/middleware"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", middleware.Timer(controller.Pong))
	fmt.Println("[GET]\t/ping")
	mux.HandleFunc("/notice/", middleware.Timer(controller.NoticeMessage))
	fmt.Println("[POST]\t/notice/")

	return mux
}
