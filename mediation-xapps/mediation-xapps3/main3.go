package main

import (
              "encoding/json"
              "log"
              "net/http"
              "github.com/gorilla/mux"
              "fmt"
              "net"
)

type rabinfo struct {
              Ecgi     string  `json:"ecgi"`
              RabId   string  `json:"rabId"`
              Qci  string  `json:"qci"`
}

// Init service and subscription var as a slice struct
var rabinfos []rabinfo


func getServices(w http.ResponseWriter, r *http.Request) {
              fmt.Println("This is TestApp3 and my IP is",  GetLocalIP())

              w.Header().Set("Content-Type", "application/json")
              json.NewEncoder(w).Encode(rabinfos)
}

// GetLocalIP returns the non loopback local IP of the host
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


// Main function
func main() {
              // Init router

              r := mux.NewRouter()

              // Hardcoded data - @todo: add database
              rabinfos = append(rabinfos, rabinfo{Ecgi: "0xBBBCCCFFFFF44", RabId: "400", Qci: "41"})
              rabinfos = append(rabinfos, rabinfo{Ecgi: "0xBBBCCCFFFFF44", RabId: "450", Qci: "49"})

              // Route handles & endpoints
              r.HandleFunc("/rnis/rab_info3", getServices).Methods("GET")
              // Start server
              log.Fatal(http.ListenAndServe(":9003", r))
}


