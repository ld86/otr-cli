package main

import  (
            "flag"
            "io"
            "net"
            "os"
            "crypto/rand"

            "golang.org/x/crypto/otr"
        )

var (
    isServer = flag.Bool("l", false, "Enable listen mode")
)

type OTRWrapper struct {
    conn io.ReadWriter
    conversation *otr.Conversation
    privateKey *otr.PrivateKey
}

func (wrapper *OTRWrapper) Read(p []byte) (n int, err error) {
    return wrapper.conn.Read(p)
}

func (wrapper *OTRWrapper) Write(p []byte) (n int, err error) {
    return wrapper.conn.Write(p)
}

func GetWrapper(conn io.ReadWriter) *OTRWrapper {
    wrapper := new(OTRWrapper)
    wrapper.conn = conn
    wrapper.conversation = new(otr.Conversation)
    wrapper.privateKey = new(otr.PrivateKey)
    wrapper.privateKey.Generate(rand.Reader)
    wrapper.conversation.PrivateKey = wrapper.privateKey;
    return wrapper
}

func chat(conn io.ReadWriter, terminal io.ReadWriter) {
    wrapper := GetWrapper(conn)

    go io.Copy(wrapper, terminal)
    io.Copy(terminal, wrapper)
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
