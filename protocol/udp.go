package protocol

import (
    "fmt"
    "io"
    "log"
    "net"
    "os"
)

func sendMulticast(message string) {
    conn, err := net.Dial("udp", MulticastAddress)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    fmt.Fprintln(conn, message)
}

func sendUnicast(ip net.Addr, message string) {
    conn, err := net.Dial("udp", ip.String()+UnicastListenAddress)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    go func() {
        mustCopy(os.Stdout, conn)
    }()
    mustCopy(conn, os.Stdin)
}

func mustCopy(dst io.Writer, src io.Reader) {
    if _, err := io.Copy(dst, src); err != nil {
        log.Fatal(err)
    }
}
