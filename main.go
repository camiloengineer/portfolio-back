package main

import (
	"net/http"

	"github.com/camiloengineer/portfolio-back/db"
	"github.com/camiloengineer/portfolio-back/routes"
	"github.com/gorilla/mux"
)

func main() {

	db.DBConnection()

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	http.ListenAndServe(":3000", r)
}
