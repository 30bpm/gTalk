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
	mux.HandleFunc("/reminder/", middleware.Timer(controller.PostReminder))
	fmt.Println("[POST]\t/reminder/")
	mux.HandleFunc("/reminder/delete/", middleware.Timer(controller.DeleteReminder))
	fmt.Println("[POST]\t/reminder/delete/")
	mux.HandleFunc("/testIP/", middleware.Timer(controller.IPTest))
	fmt.Println("[POST]\t/testIP/")
	// mux.HandleFunc("/test/", middleware.Timer(controller.TestTicker))
	// fmt.Println("[GET]\t/test/")
	return mux
}
