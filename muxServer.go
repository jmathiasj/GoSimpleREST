package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "strconv"
    "encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
)
var counter int
var mutex = &sync.Mutex{}
//Book Struct global declare
type Book struct {
	Id      string `json:"Id"`
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

func createNewBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: create New Book")
    // get the body of our POST request
    // return the string response containing the request body    
    reqBody, _ := ioutil.ReadAll(r.Body)
    var book Book
    json.Unmarshal(reqBody, &book)
    // update our global Articles array to include
    // our new Article
    Books = append(Books, book)

    json.NewEncoder(w).Encode(book)
}//creating a new book

func returnAllBooks(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: return All Books")
    json.NewEncoder(w).Encode(Books)
}//function to get all books

func returnSingleBook(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: return single Book")
    vars := mux.Vars(r)
	key := vars["id"]
	// fmt.Fprintf(w, "Key: " + key)
	for _, book := range Books {
        if book.Id == key {
            json.NewEncoder(w).Encode(book)
        }
    }

    
}//function to get single book
func updateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: update Book")
    // once again, we will need to parse the path parameters
    vars := mux.Vars(r)
    // we will need to extract the `id` of the article we
    // wish to delete
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
    var books Book
    json.Unmarshal(reqBody, &books)

    // we then need to loop through all our articles
    for index, book := range Books {
        // if our id path parameter matches one of our
        // articles
        if book.Id == id {
            // updates our Articles array to remove the 
			// article
			bookups :=Books[index+1:]
			Books = append(Books[:index],books)
			Books = append(Books,bookups...)
		}
	json.NewEncoder(w).Encode(books)
    }

}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: delete Book")
    // once again, we will need to parse the path parameters
    vars := mux.Vars(r)
    // we will need to extract the `id` of the article we
    // wish to delete
    id := vars["id"]

    // we then need to loop through all our articles
    for index, book := range Books {
        // if our id path parameter matches one of our
        // articles
        if book.Id == id {
            // updates our Articles array to remove the 
            // article
            Books = append(Books[:index], Books[index+1:]...)
        }
    }

}

func handleRequests() {
	 // creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/home", homePage)//home page


    myRouter.HandleFunc("/increment", incrementCounter)//increment counter

    myRouter.HandleFunc("/hi",func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi")
    })// direct one
    myRouter.Handle("/", http.FileServer(http.Dir("./static")))//server static folder
	myRouter.HandleFunc("/books", returnAllBooks) //get all books
	myRouter.HandleFunc("/books/{id}", returnSingleBook)//get book by id
	myRouter.HandleFunc("/book", createNewBook).Methods("POST")//post request to create a new book
	myRouter.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")//delete book
	myRouter.HandleFunc("/bookup/{id}", updateBook).Methods("PUT")
    fmt.Println("Server is running on localhost 10000")
    log.Fatal(http.ListenAndServe(":10000", myRouter))//host http server on 10000
    
}

func main() {
	fmt.Println("Rest API with MUX")
    Books = []Book{
        Book{Id: "1", Title: "Pride and Prejudice", Desc: "Novel", Author: "Jane Austem"},
        Book{Id: "2",Title: "Goosebumps", Desc: "Fiction", Author: "R.L Stine"},
    }// simulates database
    handleRequests()
}