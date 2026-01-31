package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unicode"

	"github.com/jessevdk/go-flags"
	"github.com/ochinchina/go-ini"
	apperrors "github.com/sgaunet/supervisord/internal/errors"
	"github.com/sgaunet/supervisord/internal/config"
	"github.com/sgaunet/supervisord/internal/daemon"
	"github.com/sgaunet/supervisord/internal/supervisor"
	"github.com/sgaunet/supervisord/internal/logger"
	log "github.com/sirupsen/logrus"
)

var BuildVersion string = ""

// Options the command line options
type Options struct {
	Configuration string `short:"c" long:"configuration" description:"the configuration file"`
	Daemon        bool   `short:"d" long:"daemon" description:"run as daemon"`
	EnvFile       string `long:"env-file" description:"the environment file"`
}

func init() {
	nullLogger := logger.NewNullLogger(logger.NewNullLogEventEmitter())
	log.SetOutput(nullLogger)
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		if runtime.GOOS == "windows" {
			log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})
		} else {
			log.SetFormatter(&log.TextFormatter{DisableColors: false, FullTimestamp: true})
		}
	}
	log.SetLevel(log.DebugLevel)
}

func initSignals(s *supervisor.Supervisor) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.WithFields(log.Fields{"signal": sig}).Info("receive a signal to stop all process & exit")
		s.GetManager().StopAllProcesses()
		os.Exit(-1)
	}()

}

var options Options
var parser = flags.NewParser(&options, flags.Default & ^flags.PrintErrors)

func loadEnvFile() {
	if len(options.EnvFile) == 0 {
		return
	}
	// try to open the environment file
	f, err := os.Open(options.EnvFile)
	if err != nil {
		log.WithFields(log.Fields{"file": options.EnvFile}).Error("Fail to open environment file")
		return
	}
	defer func() { _ = f.Close() }()
	reader := bufio.NewReader(f)
	for {
		// for each line
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// if line starts with '#', it is a comment line, ignore it
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] == '#' {
			continue
		}
		// if environment variable is exported with "export"
		if strings.HasPrefix(line, "export") && len(line) > len("export") && unicode.IsSpace(rune(line[len("export")])) {
			line = strings.TrimSpace(line[len("export"):])
		}
		// split the environment variable with "="
		if k, v, ok := strings.Cut(line, "="); ok {
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
			// if key and value are not empty, put it into the environment
			if len(k) > 0 && len(v) > 0 {
				_ = os.Setenv(k, v)
			}
		}
	}
}

// find the supervisord.conf in following order:
//
// 1. $CWD/supervisord.conf
// 2. $CWD/etc/supervisord.conf
// 3. /etc/supervisord.conf
// 4. /etc/supervisor/supervisord.conf (since Supervisor 3.3.0)
// 5. ../etc/supervisord.conf (Relative to the executable)
// 6. ../supervisord.conf (Relative to the executable)
func findSupervisordConf() (string, error) {
	possibleSupervisordConf := []string{options.Configuration,
		"./supervisord.ini",
		"./etc/supervisord.conf",
		"/etc/supervisord.conf",
		"/etc/supervisor/supervisord.conf",
		"../etc/supervisord.conf",
		"../supervisord.conf",
		"./supervisord.conf"}

	for _, file := range possibleSupervisordConf {
		if _, err := os.Stat(file); err == nil {
			absFile, err := filepath.Abs(file)
			if err == nil {
				return absFile, nil
			}
			return file, nil
		}
	}

	return "", apperrors.ErrConfigNotFound
}

func runServer() {
	// infinite loop for handling Restart ('reload' command)
	loadEnvFile()
	for {
		if len(options.Configuration) == 0 {
			options.Configuration, _ = findSupervisordConf()
		}
		s := supervisor.NewSupervisor(options.Configuration)
		initSignals(s)
		if _, _, _, sErr := s.Reload(true); sErr != nil {
			panic(sErr)
		}
		s.WaitForExit()
	}
}

// Get the supervisord log file
func getSupervisordLogFile(configFile string) string {
	configFileDir := filepath.Dir(configFile)
	env := config.NewStringExpression("here", configFileDir)
	myini := ini.NewIni()
	myini.LoadFile(configFile)
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	logFile := myini.GetValueWithDefault("supervisord", "logfile", filepath.Join(cwd, "supervisord.log"))
	logFile, err = env.Eval(logFile)
	if err == nil {
		return logFile
	} else {
		return filepath.Join(".", "supervisord.log")
	}
}

func main() {
	if BuildVersion != "" { supervisor.VERSION = BuildVersion }
	daemon.ReapZombie()

	// when execute `supervisord` without sub-command, it should start the server
	parser.SubcommandsOptional = true
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		if command == nil {
			log.SetOutput(os.Stdout)
			if options.Daemon {
				logFile := getSupervisordLogFile(options.Configuration)
				daemon.Daemonize(logFile, runServer)
			} else {
				runServer()
			}
			os.Exit(0)
		}
		return command.Execute(args)
	}

	if _, err := parser.Parse(); err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				_, _ = fmt.Fprintln(os.Stdout, err)
				os.Exit(0)
			default:
				_, _ = fmt.Fprintf(os.Stderr, "error when parsing command: %s\n", err)
				os.Exit(1)
			}
		}
	}
}
