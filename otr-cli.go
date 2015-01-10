package main

import  (
            "flag"
            "io"
            "net"
            "os"
            "crypto/rand"
            "fmt"

            "golang.org/x/crypto/otr"
            "golang.org/x/crypto/ssh/terminal"
        )

var (
    isServer = flag.Bool("l", false, "Enable listen mode")
    smpSecret = flag.String("s", "secret", "Secret for SMP")
)

type OTRWrapper struct {
    conn io.ReadWriter
    conversation *otr.Conversation
    privateKey *otr.PrivateKey
}

type TerminalWrapper struct {
    conn io.ReadWriter
    terminal *terminal.Terminal
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

func (wrapper *OTRWrapper) Read(p []byte) (int, error) {
    buffer := make([]byte, len(p))

    n, err := wrapper.conn.Read(buffer)
    out, _, change, toSend, _ := wrapper.conversation.Receive(buffer[:n])

    n = copy(p, out)
    SendMessages(toSend, wrapper.conn)

    switch change {
        case otr.NewKeys:
            fmt.Printf("[!] Their fingerprint: %X\n", wrapper.conversation.TheirPublicKey.Fingerprint())
            if *isServer {
                fmt.Printf("[!] Asking a question with answer: %s\n", *smpSecret)
                msgs, _ := wrapper.conversation.Authenticate("Do you know secret?", []byte(*smpSecret))
                SendMessages(msgs, wrapper.conn)
            }
        case otr.SMPSecretNeeded:
            question := wrapper.conversation.SMPQuestion()
            fmt.Printf("[!] Answer a question '%s'\n", question)
            msgs, _ := wrapper.conversation.Authenticate(question, []byte(*smpSecret))
            SendMessages(msgs, wrapper.conn)
        case otr.SMPComplete:
            fmt.Println("[!] Answer is correct")
        case otr.SMPFailed:
            fmt.Println("[!] Answer is wrong")
            os.Exit(1)
    }

    return n, err
}

func (wrapper *OTRWrapper) Write(p []byte) (int, error) {
    msgs, _ := wrapper.conversation.Send(p) // TODO(ld86) Handle all errors
    _, err := SendMessages(msgs, wrapper.conn)
    return len(string(p)), err
}

func GetOTRWrapper(conn io.ReadWriter) *OTRWrapper {
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

func (wrapper *TerminalWrapper) Read(p []byte) (int, error) {
    line, err := wrapper.terminal.ReadLine()
    line += "\n"
    return copy(p, []byte(line)), err
}

func (wrapper *TerminalWrapper) Write(p []byte) (int, error) {
    return wrapper.terminal.Write([]byte("<Their> " + string(p)))
}

func GetTerminalWrapper() *TerminalWrapper {
    wrapper := new(TerminalWrapper)
    wrapper.conn = os.Stdin
    terminal.MakeRaw(0)
    wrapper.terminal = terminal.NewTerminal(os.Stdin, "<You> ")
    return wrapper
}

func chat(conn io.ReadWriter) {
    otrWrapper := GetOTRWrapper(conn)
    terminalWrapper := GetTerminalWrapper()

    go io.Copy(terminalWrapper, otrWrapper)
    io.Copy(otrWrapper, terminalWrapper)
}

func main() {
    /*
    state, _ := terminal.MakeRaw(0)
    term := terminal.NewTerminal(os.Stdin, "$ ")
    for {
        line, err := term.ReadLine()
        term.Write([]byte(fmt.Sprintf("# %d %s\n", len(line), err)))
        if err != nil || len(line) == 0 {
            terminal.Restore(0, state)
            return
        }
        term.Write([]byte(fmt.Sprintf("$ %s\n", line)))
    }
    return
    */
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
    chat(conn)
}
