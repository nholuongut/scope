package endpoint

import (
	"github.com/nholuongut/scope/probe/endpoint/procspy"
	"github.com/nholuongut/scope/probe/process"
	"github.com/nholuongut/scope/report"
)

// Node metadata keys.
const (
	ReverseDNSNames = report.ReverseDNSNames
	SnoopedDNSNames = report.SnoopedDNSNames
	CopyOf          = report.CopyOf
)

// ReporterConfig are the config options for the endpoint reporter.
type ReporterConfig struct {
	HostID       string
	HostName     string
	SpyProcs     bool
	UseConntrack bool
	WalkProc     bool
	UseEbpfConn  bool
	ProcRoot     string
	BufferSize   int
	ProcessCache *process.CachingWalker
	Scanner      procspy.ConnectionScanner
	DNSSnooper   *DNSSnooper
}

// Name of this reporter, for metrics gathering
func (Reporter) Name() string { return "Endpoint" }
