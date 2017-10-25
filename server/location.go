package main

import (
	"github.com/oschwald/geoip2-golang"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	GeoLiteURL  = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"
	GeoLiteFile = "GeoLite2-City.mmdb"
)

func GetGeoLocation(IPAddress string) (Location local, err error) {
	_, err = os.Open(GeoLiteFile)
	var emptylocal local
	if err != nil {
		log.Print("Could not open GeoLite DB attempting to fetch over the internet...")

		err = DownloadFile(GeoLiteFile+".gz", GeoLiteURL)
		if err != nil {
			return emptylocal, err
		}
		log.Print("Completed fetching GeoLite DB")

		err = UnGzip(GeoLiteFile+".gz", ".")
		if err != nil {
			return emptylocal, err
		}
	}

	db, err := geoip2.Open(GeoLiteFile)
	if err != nil {
		return emptylocal, err
	}
	defer db.Close()

	ip := net.ParseIP(IPAddress)
	record, err := db.City(ip)
	if err != nil {
		return emptylocal, err
	}

	Location.State = record.Subdivisions[0].Names["en"]
	Location.Country = record.Country.Names["en"]
	Location.Latitude = record.Location.Latitude
	Location.Longitude = record.Location.Longitude
	return Location, nil
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