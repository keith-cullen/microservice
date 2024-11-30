package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

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

func main() {
	flags = flag.NewFlagSet("", flag.ExitOnError)
	flags.BoolVar(&opts.insecure, "i", false, "insecure (HTTP) mode")
	flags.StringVar(&opts.configFileName, "c", "", "config file name")
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "usage: %s [OPTIONS]...\n", os.Args[0])
		flags.PrintDefaults()
	}
	flags.Parse(os.Args[1:]) // ExitOnError so no need to check the return value
	if opts.insecure {
		log.Print("insecure")
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
	server, err := server.New(store)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	done := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh
		if err := server.Stop(); err != nil {
			log.Fatalf("error: %v", err)
		}
		close(done)
	}()
	err = server.Start(opts.insecure)
	if err != nil {
		log.Fatal("error: %v", err)
	}
	<-done
}
