package network

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

func Play() {
	var httpAddress = flag.String("addr", ":8080", "http address. Example :8080")
	flag.Parse()
	println("Options:\n -httpAddress:", *httpAddress)

	var wg sync.WaitGroup
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go StartSimpleHttpServer(ctx, &wg, *httpAddress)

	// wait for interupt
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	fmt.Println("Stoping...")
	cancel()
	wg.Wait()
}
