package main

import(
	"net/http"
	"fmt"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server {
	Addr:		":8080",
	Handler:	serveMux,
	}
	err := server.ListenAndServe()
	fmt.Printf("%v", err)

}
