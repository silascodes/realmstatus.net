package main

import (
    "net/http"
)

func main() {
    c := NewController()
    http.ListenAndServe(":80", c)
}
