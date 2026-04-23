package main

import (
	"flag"
	"log"

	pkg "github.com/TimoKats/roxy/pkg"
)

func main() {
	// flags
	filename := flag.String("feeds", "", "(newsboat) file with rss feeds")
	port := flag.String("port", "2112", "port number to serve on")
	flag.Parse()
	// start server
	idx := pkg.NewIndex()
	if err := idx.Load(*filename); err != nil {
		log.Println("startup failed...")
		return
	}
	idx.Serve(":" + *port)
}
