package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)
func handlePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p, _ := strconv.Atoi(vars["page"])
	w.Write(getPage(p))
}
func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/media/").Handler(http.StripPrefix("/media/",
		http.FileServer(http.Dir("./media/"))))
	r.HandleFunc("/{page}", handlePage)
	return r
}