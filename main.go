package main

import (
	"flag"
	"log"

	"github.com/jndunlap/tagit/processor"
)

func main() {
	dry := flag.Bool("dry", false, "dry-run")
	dir := flag.String("dir", ".", "root directory")
	flag.Parse()

	if err := processor.Run(*dir, *dry, log.Printf); err != nil {
		log.Fatal(err)
	}
}
