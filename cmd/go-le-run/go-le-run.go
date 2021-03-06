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
	a.Start()
	errchan := a.GetErrChan()

	//wait for a quit signal
	for a.GetShouldRun() {
		//watch for signals and errors
		select {
		//quit signal from user
		case <-sigChan:
			log.Println("Received interrupt signal")
			log.Println("Killing application...")
			go a.Stop()

		//error from app
		case appErr := <-errchan:
			if appErr != nil {
				log.Printf("App received error: %s\n", appErr)
			}
		//app should no longer run
		case <-(a.GetShouldRunChan()):
		}
	}

	log.Println("Cleaning up app")
	killed := false
	for a.GetIsRunning() && !killed {
		//watch for  signals and errors
		select {
		//second quit signal from user
		case <-sigChan:
			log.Println("Received interrupt signal")
			killed = true
		//error from app
		case appErr := <-errchan:
			if appErr != nil {
				log.Printf("App received error: %s\n", appErr)
				//a.Stop()
			}
		//app is no longer running
		case <-(a.GetIsRunningChan()):
		}
	}

	//flush errors
	outoferrs := false

	for !outoferrs {
		select {
		case err, ok := <-errchan:
			if !ok {
				outoferrs = true
			}
			if err != nil {
				log.Printf("App received error: %s\n", err)
			}
		default:
			outoferrs = true
		}

	}

	log.Println("App Finished!")
}
