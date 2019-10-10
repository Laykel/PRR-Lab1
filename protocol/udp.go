package protocol

import (
    "io"
    "log"
    "net"
    "os"
)

func sendMulticast(message string) {
    // TODO: Send bytes and not strings

    conn, err := net.Dial("udp", MulticastAddress)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    conn.Write([]byte(message))
}

func SendUnicast(message string, ip net.Addr) {
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

