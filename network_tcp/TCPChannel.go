package networktcp

import (
	"fmt"
	"log"
	"net"
)

type tcpPeer struct {
	id       uint64
	sendChan chan []byte
	recvChan chan []byte
	conn     net.Conn
	isAlive  bool
}

func newTcpPeer(id uint64, conn net.Conn) *tcpPeer {
	return &tcpPeer{
		id:       id,
		conn:     conn,
		isAlive:  true,
		sendChan: make(chan []byte),
		recvChan: make(chan []byte),
	}
}

func (p *tcpPeer) close() {
	if p.isAlive {
		close(p.sendChan)
		close(p.recvChan)

	}
	p.isAlive = false
}

func CreateTCPServer(address string) {
	var peers []*tcpPeer

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Start server. tcp:%s\n", address)

	recvMessageLoop := func(peer *tcpPeer) {
		buffer := make([]byte, 1024)
		log.Printf("conn read started %v", peer.id)

		defer func() {
			log.Println("recvMessageLoop cleanup peerId", peer.id)
			peer.close()
		}()

		for {
			_numByte, err := peer.conn.Read(buffer)
			if err != nil {
				log.Printf("rec from conn closed %v", peer.id)
				return
			}
			log.Printf("conn %v data: %d %s\n", peer.id, _numByte, string(buffer))
		}

	}

	sendMessageLoop := func(peer *tcpPeer) {
		for {
			msg := <-peer.recvChan
			_, err := peer.conn.Write(msg)
			if err != nil {
				log.Printf("send to conn closed %v", peer.id)
				return
			}
		}
	}

	onNewConnection := func(conn net.Conn) {
		id := uint64(len(peers) + 1)
		newPeer := newTcpPeer(id, conn)
		peers = append(peers, newPeer)

		defer func() {
			peers = removeIndex(peers, int(id-1))
		}()

		go sendMessageLoop(newPeer)
		recvMessageLoop(newPeer)

		newPeer.close()
	}

	accepLoop := func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Panic("server closed")
			}
			fmt.Printf("Connection Accepted %v\n", conn)
			go onNewConnection(conn)
		}
	}

	accepLoop()
}

func removeIndex[T any](s []T, index int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
