package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	logger       = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
	useEnvConfig = flag.Bool("e", false, "use environment variables as config")
)

func main() {
	flag.Usage = help
	flag.Parse()

	cmds := map[string]func(){
		"start":   startServer,
		"init":    initConfig,
		"gen-key": genKey,
		"help":    help,
	}

	if cmdFunc, ok := cmds[flag.Arg(0)]; ok {
		cmdFunc()
	} else {
		help()
		os.Exit(1)
	}
}

func help() {
	fmt.Fprintln(os.Stderr, `Usage: 
	 notes start								- start the server
	 notes init									- create an initial configuration file
	 notes gen-key							- generates a random 32-byte hex-encoded key
	 notes help 								- show this message
Use -e flag to read configuration from environment variables instead of a file. E.g.:
	 notes -e start`)
}
