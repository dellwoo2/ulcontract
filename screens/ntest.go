package main
 
import (
"fmt"
"net/http"
"html/template"
"io/ioutil"
)

func load(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/load/"):]
fmt.Print("title="+title)
str:=`<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>Go Web Programming</title>
</head>
<body>
<form action="http://127.0.0.1:8080/process?hello=world&thread=123" method="post" enctype="application/x-www-form-urlencoded">
<input type="text" name="hello" value="sau sheong"/>
<input type="text" name="post" value="456"/>
<input type="submit"/>
</form>
</body>
</html>`
fmt.Fprintln(w, str)
}
type Page struct {
    Title string
    Body  []byte
}
type Test struct{
  Title string
  Body string
  AA string
  BB string
  CC string
  DD string
  EE string
}
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    fmt.Print("Page="+title)
    fmt.Print("Page Body="+string(p.Body))
    if err != nil {
        fmt.Print(err)
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("Edit.html")
    var x Test
    x.AA="THIS IS AA"
    x.BB="THIS IS BB"
    x.CC="THIS IS CC"
    x.DD="THIS IS DD"
    x.EE="THIS IS EE"
    x.Title="THIS IS THE TITLE"
    x.Body="Some Text"
    t.Execute(w, x)
}

func loadPage(title string) (*Page ,error){
    filename := title
    body, _ := ioutil.ReadFile(filename)
    return &Page{Title: title, Body: body},nil
}
 
func process(w http.ResponseWriter, r *http.Request) {
r.ParseForm()
fmt.Fprintln(w, r.Form)
fmt.Print(r.Form)
fmt.Fprintln(w, r.FormValue("hello"))
}
 
func main() {
server := http.Server{
Addr: "127.0.0.1:8080",
}
http.HandleFunc("/process", process)
http.HandleFunc("/load/", load)
http.HandleFunc("/edit/", editHandler)
server.ListenAndServe()
}

