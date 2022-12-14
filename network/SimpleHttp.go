package network

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

type appState struct {
	atomic.Int64
}

func createAppState() *appState {
	ins := appState{}
	ins.Store(0)
	return &ins
}

func StartSimpleHttpServer(ctx context.Context, wg *sync.WaitGroup, addr string) {
	appState := createAppState()

	http.HandleFunc("/count", func(res http.ResponseWriter, req *http.Request) {

		switch req.Method {
		case "GET":
			res.WriteHeader(200)
			res.Write([]byte(fmt.Sprintf("%d", appState.Load())))
			return
		case "POST":
			appState.Add(1)
			res.WriteHeader(201)
			res.Write([]byte(fmt.Sprintf("%d", appState.Load())))
			return
		default:
			res.WriteHeader(404)
			res.Write([]byte("404 - method not support\n"))
			return
		}
	})

	log.Println("http listern at ", addr)

	httpServer := http.Server{
		Addr: addr,
	}

	wg.Add(1)
	defer wg.Done()

	go func() {
		<-ctx.Done()
		log.Println("shutdown http server begin")
		httpServer.Shutdown(ctx)
		log.Println("shutdown http server end")
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("http server start error", err)
	}
}
