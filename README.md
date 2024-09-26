# Defacto2 / magicnumber

[![Go Reference](https://pkg.go.dev/badge/github.com/Defacto2/magicnumber.svg)](https://pkg.go.dev/github.com/Defacto2/magicnumber)
[![Go Report Card](https://goreportcard.com/badge/github.com/Defacto2/magicnumber)](https://goreportcard.com/report/github.com/Defacto2/magicnumber)

The magicnumber package is to read the bytes within a file to learn of the format or type of file. This is useful for identifying the file without relying on the file extension or MIME type. See the [reference documentation](https://pkg.go.dev/github.com/Defacto2/magicnumber) for additional usage and examples.

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/Defacto2/magicnumber"
)

func main() {
    file, err := os.Open("example.exe")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    mn, err := magicnumber.Find(file)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("File type: %s\n", mn.Type)
    fmt.Printf("File extension: %s\n", mn.Extension)
}
``` 