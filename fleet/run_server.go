package fleet

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// LoadRoutes returns routes.
func LoadRoutes() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/api/robots/", newRobotListHandler()).Methods(http.MethodGet)
	r.Handle("/api/robots/{robot}/send/", newSendRobotHandler()).Methods(http.MethodPatch)
	return r
}

// RunServer runs a HTTP server.
func RunServer() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("http.port")),
		Handler: LoadRoutes(),
	}

	return server.ListenAndServe()
}
