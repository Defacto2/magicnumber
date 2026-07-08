# Defacto2 / magicnumber

[![Go Reference](https://pkg.go.dev/badge/github.com/Defacto2/magicnumber.svg)](https://pkg.go.dev/github.com/Defacto2/magicnumber)

The magicnumber package reads the bytes within a file to learn of the format or type of file, being useful for identifying the file without relying on the file extension or MIME type. See the [reference documentation](https://pkg.go.dev/github.com/Defacto2/magicnumber) for additional usage and examples.

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

    // Option 1.
    result := magicnumber.Find(file)
    fmt.Printf("File type: %s\n", result)

    // Option 2.
    valid, result, err := magicnumber.MatchExt("example.exe", file)
    if err != nil {
		    fmt.Fprintln(os.Stderr, err)
    }
    fmt.Printf("File type: %s\n", result)
    fmt.Printf("File extension valid: %v\n", valid)
}
```

