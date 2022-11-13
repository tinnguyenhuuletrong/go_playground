package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"ttin.com/play2022/network"
)

func main() {
	var wg sync.WaitGroup
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go network.StartSimpleHttpServer(ctx, &wg, ":8080")

	// wait for interupt
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	fmt.Println("Stoping...")
	cancel()
	wg.Wait()
}
