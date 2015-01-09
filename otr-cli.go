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

func SendMessages(messages [][]byte, writer io.Writer) (int, error) {
    allN := 0
    for _, msg := range messages {
        n, err := writer.Write(msg)
        allN += n
        if err != nil {
            return allN, err
        }
    }
    return allN, nil
}

func (wrapper *OTRWrapper) Read(p []byte) (n int, err error) {
    buffer := make([]byte, len(p))

    n, err = wrapper.conn.Read(buffer)
    out, _, _, msgs, _ := wrapper.conversation.Receive(buffer[:n])
    SendMessages(msgs, wrapper.conn)

    n = copy(p, out)
    return n, err
}

func (wrapper *OTRWrapper) Write(p []byte) (int, error) {
    msgs, _ := wrapper.conversation.Send(p) // TODO(ld86) Handle all errors
    return SendMessages(msgs, wrapper.conn)
}

func GetWrapper(conn io.ReadWriter) *OTRWrapper {
    wrapper := new(OTRWrapper)
    wrapper.conn = conn
    wrapper.conversation = new(otr.Conversation)
    wrapper.privateKey = new(otr.PrivateKey)
    wrapper.privateKey.Generate(rand.Reader)
    wrapper.conversation.PrivateKey = wrapper.privateKey;

    if *isServer {
        wrapper.Write([]byte(otr.QueryMessage))
    }

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
