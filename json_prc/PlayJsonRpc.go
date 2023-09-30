package play_jsonprc

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

// In a real project, these would be defined in a common file
type Args struct {
	A int
	B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func PlayJsonRpc_WSServer() {
	ins := new(Arith)
	rpc.Register(ins)

	http.Handle("/ws", websocket.Handler(func(wsc *websocket.Conn) {
		log.Printf("Handler starting")
		jsonrpc.ServeConn(wsc)
		log.Printf("Handler exiting")
	}))

	go func() {
		log.Println("Client simulation will start after 2 s;econds")
		time.Sleep(2 * time.Second)
		log.Println("Client start")
		PlayJsonRpcClient()
		log.Println("Client end")
	}()

	log.Println("ws://0.0.0.0:8000")
	err := http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		log.Panicln(err)
	}
}

func PlayJsonRpcClient() {
	origin := "http://localhost/"
	url := "ws://localhost:8000/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	args := Args{7, 8}
	var reply int

	c := jsonrpc.NewClient(ws)
	err = c.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
}

type stdInput struct {
	rw *bufio.ReadWriter
}

// Write implements io.ReadWriteCloser.
func (ins stdInput) Write(p []byte) (n int, err error) {
	n, e := ins.rw.Write(p)
	ins.rw.Flush()
	return n, e
}

// Close implements io.ReadWriteCloser.
func (ins stdInput) Close() error {
	return nil
}

// Read implements io.ReadWriteCloser.
func (ins stdInput) Read(p []byte) (n int, err error) {
	return ins.rw.Read(p)
}

var _ io.ReadWriteCloser = stdInput{}

func PlayJsonRpc_Stdio() {
	ins := new(Arith)
	rpc.Register(ins)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	readWriter := bufio.NewReadWriter(reader, writer)
	stdInp := &stdInput{rw: readWriter}

	println("Try input:", `{"id":"1","method":"Arith.Multiply","params":[{"A": 3, "B": 2}]}`)
	// input stdin {"id":"1","method":"Arith.Multiply","params":[{"A": 3, "B": 2}]}
	jsonrpc.ServeConn(stdInp)
}
