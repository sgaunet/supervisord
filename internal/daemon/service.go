package daemon

import (
	"fmt"

	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
)

// ServiceCommand install/uninstall/start/stop supervisord service.
type ServiceCommand struct {
	Configuration string
	EnvFile       string
}

type program struct{}

// Start supervised service.
func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {}

// Stop supervised service.
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func handleServiceActionResult(action string, err error) error {
	if err != nil {
		log.Errorf("Failed to %s service go-supervisord: %v", action, err)
		fmt.Printf("Failed to %s service go-supervisord: %v\n", action, err)
		return err
	}
	fmt.Printf("Succeed to %s service go-supervisord\n", action)
	return nil
}

// Execute implement Execute() method defined in flags.Commander interface, executes the given command.
func (sc ServiceCommand) Execute(args []string) error {
	if len(args) == 0 {
		showUsage()
		return nil
	}

	serviceArgs := make([]string, 0)
	if sc.Configuration != "" {
		serviceArgs = append(serviceArgs, "--configuration="+sc.Configuration)
	}
	if sc.EnvFile != "" {
		serviceArgs = append(serviceArgs, "--env-file="+sc.EnvFile)
	}

	svcConfig := &service.Config{
		Name:        "go-supervisord",
		DisplayName: "go-supervisord",
		Description: "Supervisord service in golang",
		Arguments:   serviceArgs,
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Error("service init failed", err)
		return fmt.Errorf("failed to create service: %w", err)
	}

	action := args[0]
	switch action {
	case "install":
		return handleServiceActionResult(action, s.Install())
	case "uninstall":
		_ = s.Stop()
		return handleServiceActionResult(action, s.Uninstall())
	case "start":
		return handleServiceActionResult(action, s.Start())
	case "stop":
		return handleServiceActionResult(action, s.Stop())
	default:
		showUsage()
	}

	return nil
}

func showUsage() {
	fmt.Println("usage: supervisord service install/uninstall/start/stop")
}

// RegisterServiceCommand registers the service command with the parser.
func RegisterServiceCommand(p interface {
	AddCommand(shortDescription string, longDescription string, data string, command any) (any, error)
}, serviceCmd *ServiceCommand) {
	_, _ = p.AddCommand("service",
		"install/uninstall/start/stop service",
		"install/uninstall/start/stop service",
		serviceCmd)
}
