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

type OTRWrapper struct {
}

func chat(a io.ReadWriter, b io.ReadWriter) {
    go io.Copy(a, b)
    io.Copy(b, a)
}

func main() {
    flag.Parse()
    if len(flag.Args()) != 1 {
        return
    }

    var (
        conn net.Conn
        err error
        ln net.Listener
        )

    host := flag.Args()[0]
    if *isServer {
        ln, err = net.Listen("tcp", host)
        if err != nil { return }
        conn, err = ln.Accept()
        if err != nil { return }
    } else {
        conn, err = net.Dial("tcp", host)
        if err != nil { return }
    }
    chat(conn, os.Stdin)
}
