package version

import (
	"fmt"
	"runtime"
)

var (
	VERSION  = "unknown"
	REVISION = "HEAD"
	BUILTAT  = "now"
)

func Version() string {
	version := ""
	version += fmt.Sprintf("Version:        %s\n", VERSION)
	version += fmt.Sprintf("Git hash:       %s\n", REVISION)
	version += fmt.Sprintf("Built at:       %s\n", BUILTAT)
	version += fmt.Sprintf("Golang version: %s\n", runtime.Version())
	version += fmt.Sprintf("OS/Arch:        %s/%s", runtime.GOOS, runtime.GOARCH)
	return version

}
