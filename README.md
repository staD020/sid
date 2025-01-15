# sid

sid is a Commodore 64 .sid library written in Go.

## Install Command-line Interface

`go install github.com/staD020/sid/cmd/sid@latest`

## Install Library

`go get github.com/staD020/sid`

## Use Library

```go
package main

import (
	"fmt"
	"os"
	"github.com/staD020/sid"
)

func main() {
	f, _ := os.Open("tune.sid")
	defer f.Close()
	s, _ := sid.New(f)
	fmt.Println("sid: %s", s)
	fmt.Println("bytes: %v", s.Bytes())
	return
}
```