package web

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sgaunet/supervisord/internal/supervisor"
)

// ConfAPI provides HTTP API for accessing program configuration files.
type ConfAPI struct {
	router     *mux.Router
	supervisor *supervisor.Supervisor
}

// NewConfAPI creates a ConfAPI object.
func NewConfAPI(s *supervisor.Supervisor) *ConfAPI {
	return &ConfAPI{router: mux.NewRouter(), supervisor: s}
}

// CreateHandler creates http handlers to process the program stdout and stderr through http interface.
func (ca *ConfAPI) CreateHandler() http.Handler {
	ca.router.HandleFunc("/conf/{program}", ca.getProgramConfFile).Methods("GET")
	return ca.router
}

func (ca *ConfAPI) getProgramConfFile(writer http.ResponseWriter, request *http.Request) {
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
		return nil, fmt.Errorf("failed to open config file %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}
	return b, nil
}
