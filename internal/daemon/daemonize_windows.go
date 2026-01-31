// +build windows

package daemon

func Daemonize(logfile string, proc func()) {
	proc()
}
