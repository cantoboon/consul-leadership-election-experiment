package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/lewisboon/consul-leadership-election-experiment/pkg/service"
	"log"
	"net/http"
)

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalln("Failed to create Consul client", err)
	}

	serviceManager := service.Manager{
		Services: []*service.Service{},
	}

	for i := 1; i <= 10; i++ {
		log.Println("Initialising s: ", i)
		serviceName := fmt.Sprintf("service-%d", i)
		s := service.New(serviceName, client)
		err := s.Init()
		if err != nil {
			log.Println("Failed to initialize: ", serviceName)
		}
		http.HandleFunc(fmt.Sprintf("/%d/hello", i), s.Handler)
		serviceManager.Add(s)
	}

	http.HandleFunc("/list", func(writer http.ResponseWriter, request *http.Request) {
		res := fmt.Sprintf("Total services: %d\n", len(serviceManager.Services))
		for _, s := range serviceManager.Services {
			res = res + fmt.Sprintf("%s - %s\n", s.ID, s.Status)
		}
		_, _ = writer.Write([]byte(res))
	})

	log.Println("Starting HTTP server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
