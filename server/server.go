//go:generate statik -src=../frontend

package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"github.com/rakyll/statik/fs"
	_ "github.com/alekssaul/golang_helloworld/server/statik"
)

const (
	version     = "3.1"
)

var (
	CrashAppCounter = 0
)

type EnvVars struct {
	DisplayExternalIP  bool   `default:"False"`
	DisplayGeoLocation bool   `default:"False"`
	CrashApp           bool   `default:"False"`
	CrashAppCount      int    `default:"5"`
	SimulateReady      bool   `default:"False"`
	WaitBeforeReady    int    `default:"30"`
	Port               string `default:"80"`
}

type Machine struct {
	ExternalIP  string
	LocalIP     string
	GeoLocation local
	version	string
}

type local struct{
	State string
	Country string
	Latitude float64
	Longitude float64
}

func main() {

	var env EnvVars
	var machine Machine

	err := envconfig.Process("HelloWorld", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("- Running with Flags - \nDISPLAYEXTERNALIP: %v\nDISPLAYGEOLOCATION: %v\nCrashApp: %v\nCrashAppCount: %v\nPort: %v\n",
		env.DisplayExternalIP, env.DisplayGeoLocation, env.CrashApp, env.CrashAppCount, env.Port)
	
	machine.version = version
	log.Printf("Started Application version: %s \n", machine.version)
	machine.LocalIP = GetLocalIP()
	log.Printf("Local IP: %s \n", machine.LocalIP)

	if env.DisplayExternalIP {
		machine.ExternalIP, err = GetExternalIP()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("External IP: %s \n", machine.ExternalIP)
	}

	if env.DisplayGeoLocation {
		machine.GeoLocation, err = GetGeoLocation(machine.ExternalIP)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("Location: %s , %s \n", machine.GeoLocation.State, machine.GeoLocation.Country)
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		MainApi(w, r, env, machine)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthz(w, r, env.CrashApp, env.CrashAppCount)
	})

	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		readiness(w, r, env.SimulateReady, env.WaitBeforeReady)
	})

	err = http.ListenAndServe(":"+env.Port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}
