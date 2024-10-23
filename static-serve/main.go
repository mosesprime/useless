package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
    port := flag.String("p", "3000", "local server port")
    dir := flag.String("d", "./static", "path to static asset directory")
    addr := flag.String("a", "127.0.0.1", "local server host")
    flag.Parse()

    http.Handle("/", http.FileServer(http.Dir(*dir)))

    log.Printf("Serving %s/ at %s:%s...\n", *dir, *addr, *port)
    log.Fatal(http.ListenAndServe(*addr+":"+*port, nil))
}
