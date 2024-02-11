package main

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
)

func main(){

    godotenv.Load(".env")

    fmt.Println("hello world")

    portString := os.Getenv("PORT")
    if portString == ""{
        log.Fatal("PORT is not found in the env file")
    }

    fmt.Println("PORT:", portString)
}
