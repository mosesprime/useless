package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"net"
	"time"
)

func main() {
    port := flag.String("p", "2222", "honeypot port")
    flag.Parse()

    pit := SSHTarPit{
        clients: make(map[string]sshTarClient),
        interval: 2 * time.Second,
        max_clients: 10,
    }
    fmt.Println(pit.Start(*port))
}

type SSHTarPit struct {
    clients map[string]sshTarClient
    interval time.Duration
    max_clients int
}

func (p *SSHTarPit) Start(port string) error {
    listen, err := net.Listen("tcp", "127.0.0.1:"+port)   
    if err != nil {
        return err
    }
    defer listen.Close()
    ticker := time.NewTicker(p.interval)
    go func() {
        for {
            select {
            case <-ticker.C:
                p.poll()    
            }
        }
    }()
    fmt.Printf("ssh tarpit listening on port %s\n", port)
    for {
        conn, err := listen.Accept()
        if err != nil {
            // TODO: handle accept errors
            continue
        }
        err = p.handleConn(conn)
        if err != nil {
            // TODO: handle conn error
            continue
        }
    }
}

func (p *SSHTarPit) poll() {
    chars := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ\n" // TODO: expand char set?
    var buffer bytes.Buffer
    for i := 0; i < 100; i++ { // TODO: add length variation
        n := rand.IntN(len(chars))
        buffer.WriteByte(chars[n])
    }
    payload := buffer.Bytes()
    for _, client := range p.clients {
        err := client.write(payload)
        if err != nil {
            defer client.conn.Close()
            delete(p.clients, client.remoteAddr)
            fmt.Printf("client disconnected from %s after %v\n", client.remoteAddr, time.Now().Sub(client.start_time))
        }
    }
}

func (p *SSHTarPit) handleConn(conn net.Conn) error {
    client := sshTarClient{
        conn: conn,
        start_time: time.Now(),
        remoteAddr: conn.RemoteAddr().String(),
    }
    fmt.Printf("client connected from %s\n", client.remoteAddr)
    p.clients[client.remoteAddr] = client
    if len(p.clients) >= p.max_clients {
        return errors.New("exceeded max_clients") 
    }
    return nil
}

type sshTarClient struct {
    conn net.Conn
    start_time time.Time
    remoteAddr string
}

func (c *sshTarClient) write(payload []byte) error {
    _, err := c.conn.Write(payload)
    return err
}
