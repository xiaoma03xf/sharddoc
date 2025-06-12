package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xiaoma03xf/sharddoc/lib/logger"
	"github.com/xiaoma03xf/sharddoc/tcp"
)

var cfgpath string

func init() {
	flag.StringVar(&cfgpath, "config", "", "Bootstrap node thorough this path")
}
func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Info(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()

	flag.Parse()
	if cfgpath == "" {
		fmt.Fprintf(os.Stderr, "No Bootstrap path")
		os.Exit(1)
	}
	tcp.BootstrapCluster(cfgpath)
}
