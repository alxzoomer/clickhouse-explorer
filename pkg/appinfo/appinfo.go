package appinfo

import (
	"os"
	"path/filepath"
)

var (
	Service     = filepath.Base(os.Args[0])
	Hostname, _ = os.Hostname()
	Version     = "1.0.0"
	Branch      = "dev"
	BuildTime   = "N/A"
)
