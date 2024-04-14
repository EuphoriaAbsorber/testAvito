package main

import (
	"database/sql"
	"log"
	conf "main/config"
	"net/http"
	"os"

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
	urlDB := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	log.Println(urlDB)
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

	//вспомогательные
	myRouter.HandleFunc(conf.PathFillDB, Handler.FillDB).Methods(http.MethodPost, http.MethodOptions)
	myRouter.HandleFunc(conf.PathGetUsers, Handler.GetUsers).Methods(http.MethodGet, http.MethodOptions)
	//основные
	myRouter.HandleFunc(conf.PathUserBanner, Handler.GetUserBanner).Methods(http.MethodGet, http.MethodOptions)
	myRouter.HandleFunc(conf.PathBanner, Handler.GetBanners).Methods(http.MethodGet, http.MethodOptions)
	myRouter.HandleFunc(conf.PathBanner, Handler.CreateBanner).Methods(http.MethodPost, http.MethodOptions)
	myRouter.HandleFunc(conf.PathBannerID, Handler.UpdateBanner).Methods(http.MethodPatch, http.MethodOptions)
	myRouter.HandleFunc(conf.PathBannerID, Handler.DeleteBanner).Methods(http.MethodDelete, http.MethodOptions)

	myRouter.PathPrefix(conf.PathDocs).Handler(httpSwagger.WrapHandler)

	myRouter.Use(loggingAndHeadersMiddleware)

	err = http.ListenAndServe(conf.Port, myRouter)
	if err != nil {
		log.Fatalln("cant serve")
	}
}
