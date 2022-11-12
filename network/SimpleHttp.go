package network

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func StartSimpleHttpServer(ctx context.Context, addr string) {
	var count int64 = 0
	http.HandleFunc("/count", func(res http.ResponseWriter, req *http.Request) {

		switch req.Method {
		case "GET":
			res.WriteHeader(200)
			res.Write([]byte(fmt.Sprintf("%d", count)))
			return
		case "POST":
			count += 1
			res.WriteHeader(201)
			res.Write([]byte(fmt.Sprintf("%d", count)))
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
