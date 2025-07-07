package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xiaoma03xf/sharddoc/lib/logger"
	"github.com/xiaoma03xf/sharddoc/lib/utils"
	"github.com/xiaoma03xf/sharddoc/raft"
)

var cfgpath string
var nodeID string

func init() {
	flag.StringVar(&cfgpath, "config", "", "Bootstrap node thorough this path")
	flag.StringVar(&nodeID, "node", "", "Node ID to start")
	// log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetOutput(&utils.InterceptWriter{
		W:     os.Stderr,
		Block: "Rollback failed: tx closed",
	})
}
func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Info(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()

	flag.Parse()
	if cfgpath == "" || nodeID == "" {
		fmt.Fprintf(os.Stderr, "Usage: -config <path> -node <node_id>\n")
		os.Exit(1)
	}
	raft.BootstrapCluster(cfgpath, nodeID)
}
