package main

import (
    //"log"
    "fmt"
    "os"
    //"time"

    "github.com/arachnist/gorepost/config"
    //"github.com/sorcix/irc"
)

func main() {
    config, err := config.ReadConfig(os.Args[1])
    if err != nil {
        fmt.Println("Error reading configuration from", os.Args[1], "error:", err.Error())
        os.Exit(1)
    }

    fmt.Println("My nickname:", config.Nick)
}
