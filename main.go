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

func loggingAndHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, r.Method)
		for header := range conf.Headers {
			w.Header().Set(header, conf.Headers[header])
		}

		next.ServeHTTP(w, r)
	})
}

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

	myRouter.HandleFunc(conf.PathFillDB, Handler.FillDB).Methods(http.MethodPost, http.MethodOptions)
	myRouter.HandleFunc(conf.PathGetUsers, Handler.GetUsers).Methods(http.MethodGet, http.MethodOptions)
	myRouter.HandleFunc(conf.PathUserBanner, Handler.GetUserBanner).Methods(http.MethodGet, http.MethodOptions)

	myRouter.PathPrefix(conf.PathDocs).Handler(httpSwagger.WrapHandler)

	myRouter.Use(loggingAndHeadersMiddleware)

	err = http.ListenAndServe(conf.Port, myRouter)
	if err != nil {
		log.Fatalln("cant serve")
	}
}
