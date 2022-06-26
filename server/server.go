package server

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/viveknathani/trackbit/extractor"
)

type Server struct {
	upgrader              *websocket.Upgrader
	router                *mux.Router
	shouldRunInBackground bool
	path                  string
	Storage               sync.Map
	isloadedFromDisk      bool
}

// NewServer
// Returns a new instance of Server.
func NewServer(shouldRunInBackground bool, path string) *Server {

	return &Server{
		upgrader:              &websocket.Upgrader{},
		router:                mux.NewRouter(),
		shouldRunInBackground: shouldRunInBackground,
		path:                  path,
		Storage:               *new(sync.Map),
		isloadedFromDisk:      false,
	}
}

// ServerHTTP
// Replaces the default ServeHTTP implementation.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.router.ServeHTTP(w, r)
}

// Job
// Collects stats and stores them in server's storage.
func (server *Server) Job(iface extractor.NetworkInterfaceInformation) {

	if !server.shouldRunInBackground {
		return
	}

	if !server.isloadedFromDisk {
		server.loadMapFromDisk(server.path)
		server.isloadedFromDisk = true
	}

	results := getStats(iface)
	arr, ok := server.Storage.Load(iface.Name)
	if !ok {
		arr = make([]extractor.ExtractedResults, 0)
	}
	server.Storage.Store(iface.Name, append(arr.([]extractor.ExtractedResults), *results))
	time.Sleep(time.Second)
}

// Dump
// Dumps server's in-memory storage to a trackbit.dat file under "path".
func (server *Server) Dump() {

	if !server.shouldRunInBackground {
		return
	}
	server.saveMapToDisk(server.path)
}

// SetupRoutes
// Should be called before sending server to http.ListenAndServe.
func (server *Server) SetupRoutes() {

	server.router.HandleFunc("/stats", server.pushStats)
	server.router.HandleFunc("/getHistory", server.getHistory)
	server.serveClient("client")
}

// Serve files to client.
func (server *Server) serveClient(directory string) {

	// serve index.html
	fileServer := http.FileServer(http.Dir(directory))
	server.router.Handle("/", fileServer)

	// handle paths that begin with  "/client"
	fileServer = http.StripPrefix("/"+directory, fileServer)
	server.router.PathPrefix("/" + directory + "/").Handler(fileServer)
}

// Handle the websocket connection.
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

		var results *extractor.ExtractedResults

		if server.shouldRunInBackground {
			value, ok := server.Storage.Load(networkInterfaceInfo.Name)
			if !ok {
				log.Fatal("could not find anything in server's storage")
			}
			arr := value.([]extractor.ExtractedResults)
			results = &arr[len(arr)-1]
		} else {
			results = getStats(networkInterfaceInfo)
		}

		err = connection.WriteJSON(&results)
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(time.Second)
	}
}

// Get stats from extractor package.
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

// Handle the /getHistory endpoint.
func (server *Server) getHistory(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	list, ok := server.Storage.Load(params["iface"][0])
	if !ok {
		log.Fatal("no history for this interface")
	}

	data, err := json.Marshal(list)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, ok := w.Write(data); ok != nil {
		log.Fatal(ok.Error())
	}
}
