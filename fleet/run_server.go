package fleet

import (
	"net/http"

	"github.com/gorilla/mux"
)

func loadRoutes() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/api/robots/", newRobotListHandler()).Methods(http.MethodGet)
	r.Handle("/api/robots/{robot}/send/", newSendRobotHandler()).Methods(http.MethodPatch)
	return r
}

// RunServer runs a HTTP server.
func RunServer() error {
	server := &http.Server{
		Addr:    ":8000",
		Handler: loadRoutes(),
	}

	return server.ListenAndServe()
}
