package supervisor

import (
	"fmt"
)

// VERSION the version of supervisor.

var (
	VERSION = "v0.7.3"
	COMMIT  = ""
)


// VersionCommand implement the flags.Commander interface.
type VersionCommand struct {
}

var versionCommand VersionCommand

// Execute implement Execute() method defined in flags.Commander interface, executes the given command.
func (v VersionCommand) Execute(args []string) error {
	fmt.Println("Version:", VERSION)
	fmt.Println(" Commit:", COMMIT)
	return nil
}

// RegisterVersionCommand registers the version command with the parser.
func RegisterVersionCommand(parser interface {
	AddCommand(shortDescription string, longDescription string, data string, command any) (any, error)
}) {
	_, _ = parser.AddCommand("version",
		"show the version of supervisor",
		"display the supervisor version",
		&versionCommand)
}
