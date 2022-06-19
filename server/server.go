package server

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/viveknathani/trackbit/extractor"
)

type Server struct {
	upgrader *websocket.Upgrader
	router   *mux.Router
}

func NewServer() *Server {

	return &Server{
		upgrader: &websocket.Upgrader{},
		router:   mux.NewRouter(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.router.ServeHTTP(w, r)
}

func (server *Server) SetupRoutes() {

	server.router.HandleFunc("/stats", server.pushStats)
	server.serveClient("client")
}

func (server *Server) serveClient(directory string) {

	// serve index.html
	fileServer := http.FileServer(http.Dir(directory))
	server.router.Handle("/", fileServer)

	// handle paths that begin with  "/client"
	fileServer = http.StripPrefix("/"+directory, fileServer)
	server.router.PathPrefix("/" + directory + "/").Handler(fileServer)
}

func (server *Server) pushStats(w http.ResponseWriter, r *http.Request) {

	connection, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("pushStats, upgrading connection to ws")
		log.Println(err)
		return
	}
	defer connection.Close()

	var networkInterfaceInfo extractor.NetworkInterfaceInformation
	err = connection.ReadJSON(&networkInterfaceInfo)
	if err != nil {
		log.Println("pushStats, decoding json")
		log.Println(err)
		return
	}

	for {
		results := getStats(networkInterfaceInfo)

		err = connection.WriteJSON(&results)
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(time.Second)
	}
}

func getStats(iface extractor.NetworkInterfaceInformation) *extractor.ExtractedResults {

	switch runtime.GOOS {
	case "linux":
		results := extractor.ExtractFromLinux(iface)
		return results
	case "windows":
		results := extractor.ExtractFromWindows(iface)
		return results
	}
	return nil
}
