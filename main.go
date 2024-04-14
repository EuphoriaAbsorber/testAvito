package main

import (
	"database/sql"
	"log"
	conf "main/config"
	"net/http"

	"main/delivery"
	"main/repository"
	"main/usecase"

	_ "main/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	myRouter := mux.NewRouter()
	urlDB := "postgres://" + conf.DBSPuser + ":" + conf.DBPassword + "@" + conf.DBHost + ":" + conf.DBPort + "/" + conf.DBName

	db, err := sql.Open("pgx", urlDB)
	if err != nil {
		log.Fatalln("could not connect to database")
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln("unable to reach database ", err)
	}
	log.Println("database is reachable")

	Store := repository.NewStore(db)

	Usecase := usecase.NewUsecase(Store)

	Handler := delivery.NewHandler(Usecase)

	myRouter.HandleFunc(conf.PathUserBanner, Handler.GetUserBanner).Methods(http.MethodGet, http.MethodOptions)

	myRouter.PathPrefix(conf.PathDocs).Handler(httpSwagger.WrapHandler)

	err = http.ListenAndServe(conf.Port, myRouter)
	if err != nil {
		log.Fatalln("cant serve")
	}
}
