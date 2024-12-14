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
    interval := flag.Uint("i", 2000, "interval in milliseconds") // TODO: add 0 interval = no data option
    max_clients := flag.Uint("m", 20, "maxium number of clients") // TODO: add 0 max_clients = no max option
    
    flag.Parse()

    if *interval == 0 {
        panic("interval can not be zero")
    }
    if *max_clients == 0 {
        panic("invalid max_clients")
    }

    pit := sshTarPit{
        clients: make(map[string]sshTarClient),
        interval: time.Duration(*interval) * time.Millisecond,
        max_clients: *max_clients,
    }
    fmt.Println(pit.start(*port))
}

type sshTarPit struct {
    clients map[string]sshTarClient
    interval time.Duration
    max_clients uint
}

func (p *sshTarPit) start(port string) error {
    listen, err := net.Listen("tcp", "0.0.0.0:"+port)   
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

func (p *sshTarPit) poll() {
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

func (p *sshTarPit) handleConn(conn net.Conn) error {
    client := sshTarClient{
        conn: conn,
        start_time: time.Now(),
        remoteAddr: conn.RemoteAddr().String(),
    }
    fmt.Printf("client connected from %s\n", client.remoteAddr)
    p.clients[client.remoteAddr] = client
    if uint(len(p.clients)) >= p.max_clients {
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
