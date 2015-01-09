package main

import  (
            "flag"
            "io"
            "net"
            "os"
        )

var (
    isServer = flag.Bool("l", false, "Enable listen mode")
)

func chat(a io.ReadWriter, b io.ReadWriter) {
    go io.Copy(a, b)
    io.Copy(b, a)
}

func main() {
    flag.Parse()
    if len(flag.Args()) != 1 {
        return
    }

    host := flag.Args()[0]
    if *isServer {
        ln, err := net.Listen("tcp", host)
        if err != nil {
            return
        }
        for {
            conn, err := ln.Accept()
            if err != nil { }
            chat(conn, os.Stdin)
        }
    } else {
        conn, err := net.Dial("tcp", host)
        if err != nil {
            return
        }
        chat(conn, os.Stdin)
    }
}
