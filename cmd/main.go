package main

import "fmt"

func main() {
    a := make(map[string]string)
    v, ok := a["foo"]
    if v == "" {
        fmt.Println("hello")
    }
    fmt.Println(v, ok)
}
