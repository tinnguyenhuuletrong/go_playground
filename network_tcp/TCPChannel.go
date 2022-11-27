package networktcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

// ------------
//	TcpPeer
// ------------

type TcpPeer struct {
	Id       uint64
	SendChan chan []byte
	RecvChan chan []byte
	conn     net.Conn
	isAlive  bool
}

func newTcpPeer(id uint64, conn net.Conn) *TcpPeer {
	return &TcpPeer{
		Id:       id,
		conn:     conn,
		isAlive:  true,
		SendChan: make(chan []byte, 100),
		RecvChan: make(chan []byte, 100),
	}
}

func (p *TcpPeer) String() string {
	return fmt.Sprintf("%d", p.Id)
}

func (p *TcpPeer) Send(data []byte) {
	p.SendChan <- data
}

func (p *TcpPeer) sendMessageLoop() {
	for {
		msg := <-p.SendChan
		_, err := p.conn.Write(msg)
		if err != nil {
			log.Printf("send to conn closed %v", p.Id)
			return
		}
	}
}

func (p *TcpPeer) recvMessageLoop() {
	reader := bufio.NewReader(p.conn)

	defer func() {
		log.Println("recvMessageLoop cleanup peerId", p.Id)
	}()

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("rec from conn closed %v", p.Id)
			return
		}
		p.RecvChan <- []byte(msg)
	}
}

func (p *TcpPeer) close() {
	if p.isAlive {
		close(p.SendChan)
		close(p.RecvChan)
	}
	p.isAlive = false
}

// ------------
//	TcpChannelServer
// ------------

type TcpChannelServer struct {
	mux    sync.Mutex
	nextId uint64
	peers  []*TcpPeer

	OnNewPeer chan *TcpPeer
}

func CreateTCPServer() *TcpChannelServer {
	return &TcpChannelServer{
		nextId:    0,
		peers:     make([]*TcpPeer, 0),
		OnNewPeer: make(chan *TcpPeer),
	}
}

func (s *TcpChannelServer) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Server listern at tcp:%s\n", address)
	s.acceptLoop(listener)
}

func (s *TcpChannelServer) acceptLoop(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panic("server closed")
		}
		fmt.Printf("Connection Accepted %v\n", conn)
		go s.onNewConnection(conn)
	}
}

func (s *TcpChannelServer) onNewConnection(conn net.Conn) {
	newPeer := addPeer(s, conn)

	defer func() {
		removePeer(s, newPeer)
	}()

	//TODO: Do Auth
	go func() {
		s.OnNewPeer <- newPeer
	}()

	go newPeer.sendMessageLoop()
	newPeer.recvMessageLoop()
	newPeer.close()
}

func removePeer(s *TcpChannelServer, newPeer *TcpPeer) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Remove item from slice
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i] == newPeer {
			s.peers[i] = s.peers[len(s.peers)-1]
			s.peers = s.peers[:len(s.peers)-1]
			break
		}
	}
}

func addPeer(s *TcpChannelServer, conn net.Conn) *TcpPeer {
	s.mux.Lock()
	defer s.mux.Unlock()

	id := s.nextId + 1
	s.nextId += 1
	newPeer := newTcpPeer(id, conn)
	s.peers = append(s.peers, newPeer)
	return newPeer
}
