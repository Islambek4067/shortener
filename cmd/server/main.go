
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/default", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "default")
	})

	log.Println("сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
