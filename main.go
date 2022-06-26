package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/viveknathani/trackbit/extractor"
	"github.com/viveknathani/trackbit/server"
)

func main() {

	// Parse flags
	port := flag.String("port", "8080", "Port for your server. Default is 8080")
	storagePath := flag.String("path", "", "Path to the directory that will store your history")
	flag.Parse()
	shouldRunInBackground := (*storagePath != "")

	// Manage goroutines and server shutdown signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Setup the server
	service := server.NewServer(shouldRunInBackground, *storagePath)
	service.SetupRoutes()

	// Serve
	go func() {
		err := http.ListenAndServe(":"+*port, service)
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("server started...")

	// Run in background
	if shouldRunInBackground {
		interfaceOrAdapterName := ""
		fmt.Print("Inteface or adapter name: ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		interfaceOrAdapterName = strings.TrimSuffix(line, "\n")
		interfaceOrAdapterName = strings.TrimSuffix(interfaceOrAdapterName, "\r")
		interfaceOrAdapterName = strings.TrimSuffix(interfaceOrAdapterName, "\r\n")
		fmt.Println("name is:", interfaceOrAdapterName)
		go func() {
			for {
				service.Job(extractor.NetworkInterfaceInformation{
					Name: interfaceOrAdapterName,
				})
			}
		}()
		fmt.Println("collecting info in background...")
	}

	// Terminate
	<-done
	if shouldRunInBackground {
		service.Dump()
	}
	fmt.Println("goodbye!")
}
