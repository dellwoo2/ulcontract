package main
import (
      "encoding/json"
 //    "unsafe"
    "fmt"
//     "net"
//     "net/rpc"
	"log" 
	"net/http" 
//"bytes"
 )

type Prod struct{
  result string
  Message string
  field string
}


func main() {

server := http.Server{
Addr: "127.0.0.1:81",
}

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var u Prod
        if r.Body == nil {
            http.Error(w, "Please send a request body", 400)
            return
        }
        err := json.NewDecoder(r.Body).Decode(&u)
        if err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
        fmt.Println(u.Message)
        u.Message="XXXXXXXXXXXX"
        //b := new(bytes.Buffer)
        json.NewEncoder(w).Encode(u)
    })
    log.Fatal(http.ListenAndServe(":81", nil))
    server.ListenAndServe()
}

