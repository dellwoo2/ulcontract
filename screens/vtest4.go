package main
import (
 //    "unsafe"
    "fmt"
//     "net"
//     "net/rpc"
//	"log" 
	"net/http" 
      "encoding/json"
	"io/ioutil"
//	"os"
"bytes"
"time"
"strconv"
 )
type Calc struct{
  Service string
  DOB string
  CalcDate string
  Smoker string
  Gender string
  Suminsured string
}


type Prod struct{
  Result string
  Message string
  Field string
  F1  string
}
type Res struct{
  COI string
  FMC string
  AMC string
}


func main() {
    //u := Prod{ "AAAA","BBB","CCCCC","DDDD"}
    //t := time.Now()
    //fmt.Println(strconv.Itoa(t.Day())+"/"+strconv.Itoa(t.Month())+"/"+strconv.Itoa(t.Year()) )
	year, month, day := time.Now().Date()
fmt.Println(strconv.Itoa(day)+"-"+month.String()+"-"+strconv.Itoa(year) )
    var x Calc
    x.DOB="09/08/1960"
    x.CalcDate="09/05/2017"
    x.Gender="M"
    x.Smoker="N"
    x.Suminsured="67888"
    x.Service="DemoCharges"
    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(x)
       fmt.Println(b) 
//    res, err := http.Post("http://175.141.142.92:8080/test", "application/json; charset=utf-8",  b )
   res, err := http.Post("http://192.168.0.3:8080/test", "application/json; charset=utf-8",  b )
//    res, err := http.Post("http://localhost:80/test", "application/json; charset=utf-8",  b )
    fmt.Println(err)
    //io.Copy(os.Stdout, res.Body)
    body, _ := ioutil.ReadAll(res.Body)
    var resx Res;
    json.Unmarshal(body , &resx)
    //fmt.Println(string(body))
    fmt.Println("COI="+resx.COI)
    fmt.Println("FMC="+resx.FMC)
    fmt.Println("AMC="+resx.AMC)
}


