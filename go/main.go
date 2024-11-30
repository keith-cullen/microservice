package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/server"
	"github.com/keith-cullen/microservice/store"
)

type cmdLineOpts struct {
	insecure       bool
	configFileName string
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
	flags.BoolVar(&opts.insecure, "i", false, "insecure (HTTP) mode")
	flags.StringVar(&opts.configFileName, "c", "", "config file name")
	flags.Usage = usage
	flags.Parse(os.Args[1:]) // ExitOnError so no need to check the return value
	if opts.insecure {
		log.Print("Insecure")
	}
	err := config.Open(opts.configFileName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	store, err := store.Open()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer store.Close()
	err = server.Run(store, opts.insecure)
	if err != nil {
		log.Fatal("error: %v", err)
	}
}
