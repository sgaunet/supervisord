package web

import (
	"github.com/sgaunet/supervisord/internal/supervisor"
	"net/http"

	"github.com/gorilla/mux"
)

// SupervisorWebgui the interface to show a WEBGUI to control the supervisor.
type SupervisorWebgui struct {
	router     *mux.Router
	supervisor *supervisor.Supervisor
}

// NewSupervisorWebgui create a new SupervisorWebgui object.
func NewSupervisorWebgui(s *supervisor.Supervisor) *SupervisorWebgui {
	router := mux.NewRouter()
	return &SupervisorWebgui{router: router, supervisor: s}
}

// CreateHandler create a http handler to process the request from WEBGUI.
func (sw *SupervisorWebgui) CreateHandler() http.Handler {
	sw.router.PathPrefix("/").Handler(http.FileServer(HTTP))
	return sw.router
}
