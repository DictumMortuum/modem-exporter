package main

import (
	"fmt"
	"github.com/DictumMortuum/modem-exporter/config"
	"github.com/DictumMortuum/modem-exporter/internal/metrics"
	"github.com/DictumMortuum/modem-exporter/internal/modem"
	"github.com/DictumMortuum/modem-exporter/internal/server"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
)

func main() {
	conf := config.Load()

	metrics.Init()

	serverDead := make(chan struct{})
	s := server.NewServer(conf.Port, modem.NewClient(conf))
	go func() {
		s.ListenAndServe()
		close(serverDead)
	}()

	ctx := shutdown.Context()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	select {
	case <-ctx.Done():
	case <-serverDead:
	}

	version := "0.0.5"
	fmt.Printf("modem-exporter v%s HTTP server stopped\n", version)
}
