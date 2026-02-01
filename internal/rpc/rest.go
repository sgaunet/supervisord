package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sgaunet/supervisord/internal/supervisor"

	"github.com/gorilla/mux"
	"github.com/sgaunet/supervisord/internal/types"
)

// SupervisorRestful the restful interface to control the programs defined in configuration file.
type SupervisorRestful struct {
	router     *mux.Router
	supervisor *supervisor.Supervisor
}

// NewSupervisorRestful create a new SupervisorRestful object.
func NewSupervisorRestful(s *supervisor.Supervisor) *SupervisorRestful {
	return &SupervisorRestful{router: mux.NewRouter(), supervisor: s}
}

// CreateProgramHandler create http handler to process program related restful request.
func (sr *SupervisorRestful) CreateProgramHandler() http.Handler {
	sr.router.HandleFunc("/program/list", sr.ListProgram).Methods("GET")
	sr.router.HandleFunc("/program/start/{name}", sr.StartProgram).Methods("POST", "PUT")
	sr.router.HandleFunc("/program/stop/{name}", sr.StopProgram).Methods("POST", "PUT")
	sr.router.HandleFunc("/program/log/{name}/stdout", sr.ReadStdoutLog).Methods("GET")
	sr.router.HandleFunc("/program/startPrograms", sr.StartPrograms).Methods("POST", "PUT")
	sr.router.HandleFunc("/program/stopPrograms", sr.StopPrograms).Methods("POST", "PUT")
	return sr.router
}

// CreateSupervisorHandler create http rest interface to control supervisor itself.
func (sr *SupervisorRestful) CreateSupervisorHandler() http.Handler {
	sr.router.HandleFunc("/supervisor/shutdown", sr.Shutdown).Methods("PUT", "POST")
	sr.router.HandleFunc("/supervisor/reload", sr.Reload).Methods("PUT", "POST")
	return sr.router
}

// ListProgram list the status of all the programs.
//
// json array to present the status of all programs.
func (sr *SupervisorRestful) ListProgram(w http.ResponseWriter, req *http.Request) {
	result := struct{ AllProcessInfo []types.ProcessInfo }{make([]types.ProcessInfo, 0)}
	if sr.supervisor.GetAllProcessInfo(nil, nil, &result) == nil {
		if err := json.NewEncoder(w).Encode(result.AllProcessInfo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		r := map[string]bool{"success": false}
		if err := json.NewEncoder(w).Encode(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// StartProgram start the given program through restful interface.
func (sr *SupervisorRestful) StartProgram(w http.ResponseWriter, req *http.Request) {
	defer func() { _ = req.Body.Close() }()
	params := mux.Vars(req)
	success, err := sr._startProgram(params["name"])
	r := map[string]bool{"success": err == nil && success}
	if err := json.NewEncoder(w).Encode(&r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sr *SupervisorRestful) _startProgram(program string) (bool, error) {
	startArgs := supervisor.StartProcessArgs{Name: program, Wait: true}
	result := struct{ Success bool }{false}
	err := sr.supervisor.StartProcess(nil, &startArgs, &result)
	if err != nil {
		return result.Success, fmt.Errorf("failed to start program %s: %w", program, err)
	}
	return result.Success, nil
}

// StartPrograms start one or more programs through restful interface.
func (sr *SupervisorRestful) StartPrograms(w http.ResponseWriter, req *http.Request) {
	defer func() { _ = req.Body.Close() }()
	var b []byte
	var err error

	if b, err = io.ReadAll(req.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("not a valid request"))
		return
	}

	var programs []string
	if err = json.Unmarshal(b, &programs); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("not a valid request"))
	} else {
		for _, program := range programs {
			_, _ = sr._startProgram(program)
		}
		_, _ = w.Write([]byte("Success to start the programs"))
	}
}

// StopProgram stop a program through the restful interface.
func (sr *SupervisorRestful) StopProgram(w http.ResponseWriter, req *http.Request) {
	defer func() { _ = req.Body.Close() }()

	params := mux.Vars(req)
	success, err := sr._stopProgram(params["name"])
	r := map[string]bool{"success": err == nil && success}
	if err := json.NewEncoder(w).Encode(&r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sr *SupervisorRestful) _stopProgram(programName string) (bool, error) {
	stopArgs := supervisor.StartProcessArgs{Name: programName, Wait: true}
	result := struct{ Success bool }{false}
	err := sr.supervisor.StopProcess(nil, &stopArgs, &result)
	if err != nil {
		return result.Success, fmt.Errorf("failed to stop program %s: %w", programName, err)
	}
	return result.Success, nil
}

// StopPrograms stop programs through the restful interface.
func (sr *SupervisorRestful) StopPrograms(w http.ResponseWriter, req *http.Request) {
	defer func() { _ = req.Body.Close() }()

	var programs []string
	var b []byte
	var err error
	if b, err = io.ReadAll(req.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("not a valid request"))
		return
	}

	if err := json.Unmarshal(b, &programs); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("not a valid request"))
	} else {
		for _, program := range programs {
			_, _ = sr._stopProgram(program)
		}
		_, _ = w.Write([]byte("Success to stop the programs"))
	}
}

// ReadStdoutLog read the stdout of given program.
func (sr *SupervisorRestful) ReadStdoutLog(w http.ResponseWriter, req *http.Request) {
}

// Shutdown the supervisor itself.
func (sr *SupervisorRestful) Shutdown(w http.ResponseWriter, req *http.Request) {
	defer func() { _ = req.Body.Close() }()

	reply := struct{ Ret bool }{false}
	_ = sr.supervisor.Shutdown(nil, nil, &reply)
	_, _ = w.Write([]byte("Shutdown..."))
}

// Reload the supervisor configuration file through rest interface.
func (sr *SupervisorRestful) Reload(w http.ResponseWriter, req *http.Request) {
	defer func() { _ = req.Body.Close() }()

	//nolint:dogsled // We only need the error return value
	_, _, _, err := sr.supervisor.Reload(false)
	r := map[string]bool{"success": err == nil}
	if err := json.NewEncoder(w).Encode(&r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
