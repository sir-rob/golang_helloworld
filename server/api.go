package main

import (
	"os"
	"log"
	"fmt"
	"net/http"
	"time"
	"encoding/json"
)


func healthz(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
}

func KillServer(w http.ResponseWriter, r *http.Request, CrashAppCount int) {
	if CrashAppCounter == CrashAppCount - 1 {
		os.Exit(1)
	} else {
		CrashAppCounter = CrashAppCounter + 1
		fmt.Fprintf(w, "App will crash in %v more HTTP GETs to this endpoint\n", CrashAppCount - CrashAppCounter )
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
}


