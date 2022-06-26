package server

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/viveknathani/trackbit/extractor"
)

// takes sync.Map, converts it into map, stores it on disk
func (server *Server) saveMapToDisk(path string) {

	file, err := os.Create(path + "/trackbit.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	err = enc.Encode(fromInternalMap(&server.Storage))
	if err != nil {
		log.Fatal(err)
	}
}

// takes json from disk into map, converts it into sync.Map
func (server *Server) loadMapFromDisk(path string) {

	file, err := os.Open(path + "/trackbit.json")
	if err != nil && os.IsNotExist(err) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var temp map[string][]extractor.ExtractedResults
	dec := json.NewDecoder(file)
	err = dec.Decode(&temp)
	if err != nil {
		log.Fatal(err)
	}
	server.Storage = *toInternalMap(&temp)
}

func toInternalMap(temp *map[string][]extractor.ExtractedResults) *sync.Map {

	m := sync.Map{}
	for key, value := range *temp {
		m.Store(key, value)
	}
	return &m
}

func fromInternalMap(m *sync.Map) *map[string][]extractor.ExtractedResults {

	temp := make(map[string][]extractor.ExtractedResults)
	m.Range(func(key interface{}, value interface{}) bool {
		temp[key.(string)] = value.([]extractor.ExtractedResults)
		return true
	})
	return &temp
}
