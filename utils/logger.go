package utils

import (
	stdlog "log"
	"os"
)

var Log stdlog.Logger = *stdlog.New(os.Stdout, "", stdlog.Lshortfile|stdlog.LUTC)
