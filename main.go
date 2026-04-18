package main

import (
	"flag"
	"log"

	pkg "github.com/TimoKats/roxy/pkg"
)

func main() {
	filename := flag.String("filename", "", "(newsboat) file with rss feeds")
	port := flag.String("port", "2112", "port number to serve on")
	log.Println("starting roxy...")
	idx := pkg.NewIndex()
	idx.Load(*filename)
	idx.Serve(":" + *port)
}
