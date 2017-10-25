package main

import (
	"log"
	"fmt"
	"net/http"
	"time"
    "encoding/json"
)


func healthz(w http.ResponseWriter, r *http.Request, CrashApp bool, CrashAppCount int) {
	if CrashApp && CrashAppCounter >= CrashAppCount {
		// do nothing
	} else {
		w.Write([]byte("OK"))
	}
}

func readiness(w http.ResponseWriter, r *http.Request, SimulateReady bool, WaitBeforeReady int) {
	if SimulateReady {
		time.Sleep(time.Duration(WaitBeforeReady) * time.Second)
	}
	w.Write([]byte("OK"))
}

func MainApi(w http.ResponseWriter, r *http.Request, env EnvVars, machine Machine) {
	b, err := json.Marshal(machine)
    if err != nil {
        log.Fatal(err)
    }
	
	fmt.Fprintf(w, "%s\n", b)

	CrashAppCounter = CrashAppCounter + 1
}


