package main

import (
	"compress/gzip"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/oschwald/geoip2-golang"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	version     = "3.1"
	GeoLiteURL  = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"
	GeoLiteFile = "GeoLite2-City.mmdb"
)

var (
	CrashAppCounter = 0
)

type EnvVars struct {
	DisplayExternalIP  bool   `default:"False"`
	DisplayGeoLocation bool   `default:"False"`
	CrashApp           bool   `default:"False"`
	CrashAppCount      int    `default:"5"`
	Debug              bool   `default:"False"`
	SimulateReady      bool   `default:"False"`
	WaitBeforeReady    int    `default:"30"`
	Port               string `default:"80"`
}

type Machine struct {
	ExternalIP  string
	LocalIP     string
	GeoLocation string
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetGeoLocation(IPAddress string) (IP string, err error) {
	_, err = os.Open(GeoLiteFile)
	if err != nil {
		log.Print("Could not open GeoLite DB attempting to fetch over the internet...")
	}

	err = downloadFile(GeoLiteFile+".gz", GeoLiteURL)
	if err != nil {
		return "", err
	}
	log.Print("Completed fetching GeoLite DB")

	err = ungzip(GeoLiteFile+".gz", ".")
	if err != nil {
		return "", err
	}

	db, err := geoip2.Open(GeoLiteFile)

	if err != nil {
		return "", err
	}
	defer db.Close()

	ip := net.ParseIP(IPAddress)
	record, err := db.City(ip)
	if err != nil {
		return "", err
	}
	LocationString := record.Subdivisions[0].Names["en"] + ", " + record.Country.Names["en"]

	return LocationString, nil
}

func GetExternalIP() (ExternalIP string, err error) {

	resp, err := http.Get("http://myexternalip.com/raw")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return strings.TrimSpace(string(body[:])), err
}

func HelloWorld(w http.ResponseWriter, r *http.Request, env EnvVars, machine Machine) {
	fmt.Fprintf(w, "Hello, world! version: %s ", version)
	fmt.Fprintf(w, "Local IP: %s \n", machine.LocalIP)
	if env.DisplayExternalIP {
		fmt.Fprintf(w, "External IP: %s \n", machine.ExternalIP)
	}
	if env.DisplayGeoLocation {
		fmt.Fprintf(w, "Location: %s \n", machine.GeoLocation)
	}
	CrashAppCounter = CrashAppCounter + 1
}

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

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ungzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

func main() {

	var env EnvVars
	var machine Machine

	err := envconfig.Process("HelloWorld", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	if env.Debug {
		fmt.Printf("DisplayExternalIP: %v\nDisplayGeoLocation: %v\nCrashApp: %v\nCrashAppCount: %v\nPort: %v\n",
			env.DisplayExternalIP, env.DisplayGeoLocation, env.CrashApp, env.CrashAppCount, env.Port)
	}

	log.Printf("Started Application version: %s \n", version)
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
		log.Printf("Location: %s \n", machine.GeoLocation)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HelloWorld(w, r, env, machine)
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
