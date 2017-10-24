package location

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

func GetGeoLocation(IPAddress string) (IP string, err error) {
	_, err = os.Open(GeoLiteFile)
	if err != nil {
		log.Print("Could not open GeoLite DB attempting to fetch over the internet...")
	}

	err = DownloadFile(GeoLiteFile+".gz", GeoLiteURL)
	if err != nil {
		return "", err
	}
	log.Print("Completed fetching GeoLite DB")

	err = UnGzip(GeoLiteFile+".gz", ".")
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
