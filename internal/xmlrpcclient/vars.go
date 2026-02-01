// Package xmlrpcclient provides XML-RPC client functionality for communicating with supervisord.
package xmlrpcclient

// https://github.com/Supervisor/supervisor/blob/ff7f18169bcc8091055f61279d0a63997d594148/supervisor/xmlrpc.py#L26-L44.
var (
	UnknownMethod       = 1
	IncorrectParameters = 2
	BadArguments        = 3
	SignatureUnsupported = 4
	ShutdownState       = 6
	BadName             = 10
	BadSignal           = 11
	NoFile              = 20
	NotExecutable       = 21
	Failed              = 30
	AbnormalTermination = 40
	SpawnError          = 50
	AlreadyStarted      = 60
	NotRunning          = 70
	Success             = 80
	AlreadyAdded        = 90
	StillRunning        = 91
	CantReread          = 92
)

// ProcStatusInfo contains status information for a single process.
type ProcStatusInfo struct {
	Name        string `xml:"name" json:"name"`
	Group       string `xml:"group" json:"group"`
	Status      int    `xml:"status" json:"status"`
	Description string `xml:"description" json:"description"`
}

// AllProcStatusInfoReply contains status information for all processes.
type AllProcStatusInfoReply struct {
	Value []ProcStatusInfo
}
