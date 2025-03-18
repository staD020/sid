package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/staD020/sid"
)

const Version = 0.2

func main() {
	flag.Parse()
	multipleFiles := flag.NArg() > 1
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		s, err := sid.New(f)
		if err != nil {
			log.Fatal(err)
		}
		if multipleFiles {
			fmt.Printf("%q: ", path)
		}
		fmt.Println(s)
	}
}
