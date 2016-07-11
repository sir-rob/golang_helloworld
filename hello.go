package main

import (
  "fmt"
  "net/http"
  "net"
  "os" 
  "io/ioutil"
  "strings"
  "github.com/oschwald/geoip2-golang"
)

const (
  port = ":80"
  version = "1.0"
  ReportExternalIP = true
  ReportGeoLocation = true
)

var (
  LocalIP = GetLocalIP()
  ExternalIP = GetExternalIP() 
  GeoLocation = GetGeoLocation(ExternalIP) 
)

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

func GetGeoLocation(IPAddress string) string {
    db, err := geoip2.Open("GeoLite2-City.mmdb")
    if err != nil {
      os.Stderr.WriteString(err.Error())
    }
    defer db.Close()
    // If you are using strings that may be invalid, check that ip is not nil
    ip := net.ParseIP(IPAddress)
    record, err := db.City(ip)
    if err != nil {
      os.Stderr.WriteString(err.Error())
    }
    LocationString := record.Subdivisions[0].Names["en"] + ", " + record.Country.Names["en"]

    return LocationString
}

func GetExternalIP() string {

  resp, err := http.Get("http://myexternalip.com/raw")

  if err != nil {
    os.Stderr.WriteString(err.Error())
    os.Exit(1)
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  return strings.TrimSpace(string(body[:]))
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {  
  fmt.Fprintf(w, "Hello, world! version: %s ", version)
  fmt.Fprintf(w, "Local IP: %s \n", LocalIP)
  if ReportExternalIP ==true { fmt.Fprintf(w, "External IP: %s \n", ExternalIP) }
  if ReportGeoLocation ==true { fmt.Fprintf(w, "Location: %s \n", GeoLocation) } 
}

func init() {
  fmt.Printf("Started Application version: %s \n", version)
  fmt.Printf("Local IP: %s \n", LocalIP)  
  if ReportExternalIP ==true { fmt.Printf("External IP: %s \n", ExternalIP) }
  if ReportGeoLocation ==true { fmt.Printf("Location: %s \n", GeoLocation) } 
  http.HandleFunc("/", HelloWorld)
  http.ListenAndServe(port, nil)
}

func main() {
  GetExternalIP()
  
}
