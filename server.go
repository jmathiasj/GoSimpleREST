package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "strconv"
    "encoding/json"
    "github.com/gorilla/mux"
)
var counter int
var mutex = &sync.Mutex{}
//Book Struct global declare
type Book struct {
    Title string `json:"Title"`
    Desc string `json:"desc"`
    Author string `json:"author"`
}
var Books []Book

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    counter++
    fmt.Fprintf(w, strconv.Itoa(counter))
    mutex.Unlock()
}//synchronization

func returnAllBooks(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: return All Books")
    json.NewEncoder(w).Encode(Books)
}

func handleRequests() {
    http.HandleFunc("/home", homePage)

    http.HandleFunc("/hey", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })

    http.HandleFunc("/increment", incrementCounter)

    http.HandleFunc("/hi",func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi")
    })// direct one
    http.Handle("/", http.FileServer(http.Dir("./static")))//server static folder
    http.HandleFunc("/books", returnAllBooks) //get all books
    fmt.Println("Server is running on localhost 10000")
    log.Fatal(http.ListenAndServe(":10000", nil))//host http server on 10000
    
}

func main() {
    Books = []Book{
        Book{Title: "Pride and Prejudice", Desc: "Novel", Author: "Jane Austem"},
        Book{Title: "Goosebumps", Desc: "Fiction", Author: "R.L Stine"},
    }// simulates database
    handleRequests()
}