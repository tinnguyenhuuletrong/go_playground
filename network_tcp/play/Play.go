package networktcp_play

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	networktcp "ttin.com/play2022/network_tcp"
)

func _server(wg *sync.WaitGroup) {
	defer wg.Done()

	server := networktcp.CreateTCPServer()
	onNewPeerCallback := func() {
		for {
			peer := <-server.OnNewPeer
			log.Println("OnNewPeer", peer.Id)
			peer.Send([]byte(fmt.Sprintf("%d\n", peer.Id)))

			func(peer *networktcp.TcpPeer) {
				for {
					bData, isClosed := <-peer.RecvChan
					if !isClosed {
						log.Printf("OnNewPeer %d closed \n", peer.Id)
						return
					}
					log.Printf("OnNewPeer %d data: %s", peer.Id, string(bData))
				}
			}(peer)
		}

	}
	go onNewPeerCallback()

	server.Start("localhost:3000")
}

func _client(wg *sync.WaitGroup) {
	defer wg.Done()

	doCreateClient := func() {
		defer wg.Done()
		tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:3000")
		if err != nil {
			log.Panic("ResolveTCPAddr failed:", err.Error())
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Panic("Dial failed:", err.Error())
		}

		reply := make([]byte, 1024)

		_, err = conn.Read(reply)
		if err != nil {
			log.Panic("Write to server failed:", err.Error())
		}

		peerId := string(reply)
		println("reply from server=", peerId)

		strEcho := "hi_from_id=" + peerId

		_, err = conn.Write([]byte(strEcho))
		if err != nil {
			log.Panic("Write to server failed:", err.Error())
		}

		println("write to server = ", strEcho)

		conn.Close()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doCreateClient()
	}
}

func PlayTcpChannel() {
	var wg sync.WaitGroup
	wg.Add(2)
	go _server(&wg)

	time.Sleep(1 * time.Second)

	go _client(&wg)

	wg.Wait()
}
