package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/keith-cullen/microservice/server"
)

type cmdLineOpts struct {
	tlsEn bool
}

var (
	opts  cmdLineOpts
	flags *flag.FlagSet
)

func usage() {
	writer := flags.Output()
	fmt.Fprintf(writer, "usage: %s [OPTIONS]...\n", os.Args[0])
	flags.PrintDefaults()
}

func main() {
	flags = flag.NewFlagSet("", flag.ExitOnError)
	flags.BoolVar(&opts.tlsEn, "s", false, "enable TLS")
	flags.Usage = usage
	flags.Parse(os.Args[1:]) // ExitOnError so no need to check the return value
	if opts.tlsEn {
		log.Print("TLS enabled")
	} else {
		log.Print("TLS not enabled")
	}
	server, err := server.Open(opts.tlsEn)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer server.Close()
	log.Printf("listening on %s", server.Address())
	log.Fatal(server.Run())
}
