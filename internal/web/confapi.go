package web

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sgaunet/supervisord/internal/supervisor"
)

type ConfApi struct {
	router     *mux.Router
	supervisor *supervisor.Supervisor
}

// NewConfApi creates a ConfApi object.
func NewConfApi(s *supervisor.Supervisor) *ConfApi {
	return &ConfApi{router: mux.NewRouter(), supervisor: s}
}

// CreateHandler creates http handlers to process the program stdout and stderr through http interface.
func (ca *ConfApi) CreateHandler() http.Handler {
	ca.router.HandleFunc("/conf/{program}", ca.getProgramConfFile).Methods("GET")
	return ca.router
}

func (ca *ConfApi) getProgramConfFile(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	if vars == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	programName := vars["program"]
	programConfigPath := getProgramConfigPath(programName, ca.supervisor)
	if programConfigPath == "" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := readFile(programConfigPath)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(b)
}

func getProgramConfigPath(programName string, s *supervisor.Supervisor) string {
	cfg := s.GetConfig()
	c := cfg.GetProgram(programName)
	if c == nil {
		return ""
	}
	res := c.GetString("conf_file", "")
	return res
}

func readFile(path string) ([]byte, error) {
	// #nosec G304 - path is validated at caller from supervisor configuration
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return io.ReadAll(f)
}
