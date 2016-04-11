package main

import (
  "fmt"
  "net/http"
)

const (
  port = ":80"
  version = "1.0"
)

var calls = 0

func HelloWorld(w http.ResponseWriter, r *http.Request) {
  calls++
  fmt.Fprintf(w, "Hello, world! You have called me %d times.\nversion: %s \n", calls, version)
}

func init() {
  fmt.Printf("Started server version:%s at http://localhost%v.\n", port, version)
  http.HandleFunc("/", HelloWorld)
  http.ListenAndServe(port, nil)
}

func main() {}
