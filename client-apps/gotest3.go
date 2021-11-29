package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "os"
   "strconv"
   "fmt"
   "time"
)

func main() {
   count, err := strconv.Atoi(os.Args[1])
   if err == nil {
      fmt.Println(count)
      }
   for {
   for i := 0; i < count; i++ {
   resp, err := http.Get("http://localhost:8000/rnis/rab_info3")
   if err != nil {
      log.Fatal(err)
   }
//We Read the response body on the line below.
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatal(err)
   }
//Convert the body to type string
   sb := string(body)
   log.Printf(sb)
}
time.Sleep(2 * time.Second)
}
}

