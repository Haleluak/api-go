package main

import (
	"github.com/Haleluak/kb-backend/app"
	"github.com/Haleluak/kb-backend/app/controller"
	"github.com/Haleluak/kb-backend/config/global"
	"github.com/Haleluak/kb-backend/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main(){
	log.Printf("Server started. Listening on port %s", global.PORT)
	log.Printf("UTC Time: %s", time.Now().UTC())
	models.InitDB()
	router := mux.NewRouter()
	router.HandleFunc("/testConnection", testConnection).Methods(global.HTTP_GET)
	router.HandleFunc("/api/user/new", controller.CreateAccount).Methods(global.HTTP_POST)
	router.HandleFunc("/api/user/login", controller.Authenticate).Methods(global.HTTP_POST)
	router.HandleFunc("/api/user/update", controller.UpdateUser).Methods(global.HTTP_POST)
	router.HandleFunc("/api/user/loginFb", controller.LoginFacebook).Methods(global.HTTP_POST)
	router.HandleFunc("/api/question/new", controller.CreateQuestion).Methods(global.HTTP_POST)
	router.HandleFunc("/api/question/{id}", controller.GetQuestion).Methods(global.HTTP_GET)
	router.HandleFunc("/api/answer/new", controller.CreateAnswer).Methods(global.HTTP_POST)
	router.HandleFunc("/api/home", controller.GetQuestions).Methods(global.HTTP_GET)

	router.Use(app.JwtAuthentication) //attach JWT auth middleware
	// Start HTTP server async
	go startHTTPServer(router)
	// Run program until interrupted.
	waitExitSignal()
}

// Start HTTP server.
func startHTTPServer(handler http.Handler) {
	err := http.ListenAndServe(global.PORT, handler)
	if err != nil {
		log.Fatal(err)
	}
}

// Wait until program interrupted.
func waitExitSignal() {
	wait := make(chan int)
	<-wait
	log.Printf("Server stop. Time: %s", time.Now().UTC())
}