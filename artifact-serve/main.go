package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
    port := flag.String("p", "3000", "local server port")
    // TODO: memorize := flag.Uint("m", 3000, "size of file (in KiB) that will get cached in memory")
    // TODO: daemonize := flag.Bool("d", false, "run in daemon mode")

    flag.Parse()
    
    artifacts := make(map[string][]byte) // TODO: add eviction & proper storage

    mux := http.NewServeMux()

    mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id")
        v, ok := artifacts[id]
        if ok {
            log.Printf("hit %s from %s\n", id, r.RemoteAddr)
            w.Write(v)
        } else {
            log.Printf("miss %s from %s\n", id, r.RemoteAddr)
            http.Error(w, "not found", http.StatusNotFound)
        }
    })

    mux.HandleFunc("PUT /", func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        defer r.Body.Close()
        hash := sha256.New()
        hash.Write(body)
        id := base64.URLEncoding.EncodeToString(hash.Sum(nil))
        log.Printf("new artifact %s from %s\n", id, r.RemoteAddr)
        fmt.Printf("\n\tCURL DOWNLOAD:\tcurl -o <DEST_FILE> <ADDR>:%s/%s\n\n", *port, id)
        artifacts[id] = body
        w.WriteHeader(http.StatusOK)
    })

    log.Printf("starting server at 127.0.0.1:%s ...\n", *port)
    
    fmt.Printf("\n\tCURL UPLOAD:\tcurl -T <FILE> 127.0.0.1:%s\n\n", *port)

    log.Fatalln(http.ListenAndServe("127.0.0.1:"+*port, mux))
}
