package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func main() {
    file, err := os.Open("../go.sum")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        parts := strings.Fields(scanner.Text())
		b64 := strings.TrimPrefix(parts[2], "h1:")
		sha256Bytes, _ := base64.StdEncoding.DecodeString(b64) //DecodeString(b64).
        if len(parts) == 3 {
            fmt.Printf("Module: %s, Version: %s, Checksum: %x\n", parts[0], parts[1], sha256Bytes)
        }
    }
}
