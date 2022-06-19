package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/viveknathani/trackbit/server"
)

func main() {

	service := server.NewServer()
	service.SetupRoutes()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	port := flag.String("port", "8080", "Port for your server. Default is 8080")
	flag.Parse()
	go func() {
		err := http.ListenAndServe(":"+*port, service)
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("server started...")
	<-done
	fmt.Println("goodbye!")
}
