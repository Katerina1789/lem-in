package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: lem-in <input-file>")
        os.Exit(1)
    }

    path := os.Args[1]
    f, err := os.Open(path)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer f.Close()

    if _, err := io.Copy(os.Stdout, f); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
