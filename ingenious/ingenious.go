package main

import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus"
    "net/http"
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "strconv"
    "time"
    "os/exec"
)

type route struct {
    route_name string
    current_count int
    prev_count int
    scale bool
    stable_ct_high int
    stable_ct_low int
}

var (
    svc_metrics = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "sep_api_processed_per_seconds",
        Help: "This contains the events processed per second by all the services",
     }, []string{"service_name", "route" })
)

var counters = make(map[string]route)
var up_threshold_value = 200
var down_threshold_value = 50

func decode_metrics(string_metrics string) {
    resp_str := strings.Split(string(string_metrics), "\n")
    for _, val := range resp_str {
        if(strings.Contains(val, "kong_http_status") && strings.Contains(val, "service=") && strings.Contains(val, "code=\"200\"")) {
            words := strings.FieldsFunc(val, func(r rune) bool { return strings.ContainsRune(" \",{}=", r) })
            svc_name, rt_name := words[2], words[4]
            count, err :=  strconv.Atoi(words[7])
            if err != nil {
                log.Fatal(err)
            }
            counters[svc_name] = route{route_name:rt_name, current_count:count, scale:false, prev_count:0, stable_ct_high:0, stable_ct_low:0}
        }
    }
}

func scale_service(app_name string, direction int) {
    app := "/usr/local/bin/kubectl"
    arg0 := "scale"
    arg1 := "deployments"
    arg2 := "--replicas=1"
    if( direction == 1) {
        arg2 = "--replicas=2"
    }
    arg3 := app_name
    arg4 := "-n"
    arg5 := "default"

    fmt.Println(app, arg0, arg1, arg2, arg3, arg4, arg5)
    cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5)
    stdout, err := cmd.Output()
    
    if err != nil {
        fmt.Println(err)
        return
    }

    // Print the output
    fmt.Println(string(stdout))
}

func read_metrics_and_orchestrate() {
    for {
        resp, err := http.Get("http://10.45.134.220:8001/metrics")
        if err != nil {
            log.Fatal(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }
        resp_str := strings.Split(string(body), "\n")
        for _, val := range resp_str {
            if(strings.Contains(val, "kong_http_status") && strings.Contains(val, "service=") && strings.Contains(val, "code=\"200\"")) {
                words := strings.FieldsFunc(val, func(r rune) bool { return strings.ContainsRune(" \",{}=", r) })
                svc_name, rt_name := words[2], words[4]
                cur_ct, err :=  strconv.Atoi(words[7])
                if err != nil {
                    log.Fatal(err)
                }
                svc_info := counters[svc_name]
                svc_info.current_count=cur_ct
                request_processed_per_second := svc_info.current_count - svc_info.prev_count
                
                // Send counters to prometheous here only...
                svc_metrics.With(prometheus.Labels{"service_name": svc_name, "route": rt_name}).Set(float64(request_processed_per_second))
                
                // Threshold, --> send scaling request.
                // fmt.Println("Debug_svc:", svc_name, " prev_ct:", svc_info.prev_count, " cur_ct:", svc_info.current_count, " scale_val:", svc_info.scale, "request_processed_per_second:", request_processed_per_second, "stable_ct_high:", svc_info.stable_ct_high, "stable_ct_low:", svc_info.stable_ct_low)

                svc_info.prev_count=svc_info.current_count
                if(request_processed_per_second!=0 && request_processed_per_second>up_threshold_value && svc_info.stable_ct_high>=3 && !svc_info.scale ) {
                    fmt.Println("Need to scale up ", svc_name, " with route: ", rt_name)
                    scale_service(svc_name, 1)
                    svc_info.stable_ct_high=0
                    svc_info.scale = true
                } else if(request_processed_per_second>up_threshold_value && !svc_info.scale) {
                    svc_info.stable_ct_high++
                } else {
                    svc_info.stable_ct_high=0
                }

                if(request_processed_per_second!=0 && request_processed_per_second<down_threshold_value && svc_info.stable_ct_low>=3 && svc_info.scale) {
                    fmt.Println("Need to scale down", svc_name, " with route: ", rt_name)
                    scale_service(svc_name, 0)
                    svc_info.stable_ct_low=0
                    svc_info.scale = false
                } else if (request_processed_per_second<down_threshold_value && svc_info.scale) {
                    svc_info.stable_ct_low++
                } else {
                    svc_info.stable_ct_low=0
                }

                counters[svc_name] = svc_info
            }
        }
        time.Sleep(2*time.Second)
    }
}

func main() {
    resp, err := http.Get("http://10.45.134.220:8001/metrics")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    decode_metrics(string(body))

    go read_metrics_and_orchestrate() 

    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":8050", nil)
}
