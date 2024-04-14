package main

import (
	"log"
	conf "main/config"
	"net/http"

	"github.com/gorilla/mux"
	//httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	myRouter := mux.NewRouter()
	err := http.ListenAndServe(conf.Port, myRouter)
	if err != nil {
		log.Println("cant serve")
	}
}
