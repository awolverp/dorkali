package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/awolverp/dorkali"
	_ "github.com/awolverp/dorkali/google"
)

const (
	Version         = "v1.1.3"
	VersionMesssage = "dorkali " + Version + " ( by awolverp ) / %s\n"
	UsageMessage    = "Dorkali a program written in golang to dorks queries in search engines\n\n" +
		"Usage:\n" +
		"\t%s [list | version [engineName] | help [engineName]]\n" +
		"\t%s engineName [OPTIONS]\n\n" +
		"*Commands:\n" +
		"\tversion [engineName]   print version, or engine version if pass engineName, and exit\n" +
		"\tlist                   print list of engines and exit\n" +
		"\thelp [engineName]      print this help, or print engine help if pass engineName, and exit\n"
)

var engine *dorkali.API = nil

func main() {

	if len(os.Args) < 2 {
		fmt.Printf(UsageMessage, os.Args[0], os.Args[0])
		return
	}

	switch os.Args[1] {
	// version
	case "version":
		if len(os.Args) == 3 {
			fmt.Println(os.Args[2] + " version: " + UseEngineOrExit(os.Args[2]).Version())
			return
		}

		fmt.Printf(VersionMesssage, runtime.Version())
		return

	// list
	case "list":
		s := dorkali.Engines()
		println("Registered engines:")
		for _, v := range s {
			println("\t" + v)
		}
		return

	// help
	case "help":
		if len(os.Args) == 3 {
			UseEngineOrExit(os.Args[2]).Usage()
			return
		}

		fmt.Printf(UsageMessage, os.Args[0], os.Args[0])
		return

	// Use engine
	default:
		engine = UseEngineOrExit(os.Args[1])
	}

	if err := engine.Start(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	response, err := engine.Search(os.Args[2:])
	if err != nil {
		fmt.Printf("error on search: %s\n", err.Error())
		os.Exit(1)
	}

	results, err := engine.ParseResponse(response)
	if err != nil {
		fmt.Printf("error on parsing: %s\n", err.Error())
		os.Exit(1)
	}

	for _, r := range results {
		fmt.Println(r.Url())
	}
}

func UseEngineOrExit(engineName string) *dorkali.API {
	e, err := dorkali.UseWithoutStart(engineName)
	if err != nil {
		fmt.Printf("Engine %q not registered! use `%s list` to see engines.\n", engineName, os.Args[0])
		os.Exit(1)
	}

	return e
}
