package main

import (
    "net/http"
)

func main() {
    http.HandleFunc("/cafe", mainHandle)
    http.ListenAndServe(":8080", nil)
}
