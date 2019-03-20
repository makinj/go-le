package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/makinj/go-le/internal/app"
	"github.com/makinj/go-le/pkg/config"
)

func main() {

	//parse flags
	cfgDir := flag.String("config-dir", "/etc/go-le/", "The directory where go-le's configuration files can be found.")
	flag.StringVar(cfgDir, "c", "/etc/go-le/", "config-dir (shortcut)")

	flag.Parse()

	//catch interrupts
	log.Printf("(Press Ctrl+c to quit)")
	var sigChan chan os.Signal
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	//load app config
	log.Printf("Loading app from: %s\n", *cfgDir)
	var cfg app.Config
	err := config.Load(*cfgDir, "app", &cfg)
	if err != nil {
		log.Fatalf("Error loading configuration for application: %s", err)
	}

	//Create app from config
	a, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Error creating application: %s", err)
	}

	//Start app
	log.Println("Starting App")
	errchan := a.Run()

	//watch for signals and errors
	select {
	case <-sigChan:
		log.Println("Received interrupt signal")
		log.Println("Killing application...")
		a.Kill()
	case appErr := <-errchan:
		if appErr != nil {
			log.Printf("App received error: %s\n", appErr)
			a.Kill()
		}
	}

	<-errchan

	log.Println("App Finished!")
}
