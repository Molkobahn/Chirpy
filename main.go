package main

import(
	"net/http"
	"fmt"
)

func main() {
	serverMux := http.NewServeMux()
	server := http.Server {
	Addr:		":8080",
	Handler:	serverMux,
	}
	err := server.ListenAndServe()
	fmt.Printf("%v", err)

}
