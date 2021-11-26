package main

import (
	"consul-leadership-election-example/pkg/service"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
)

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalln("Failed to create Consul client", err)
	}

	serviceManager := service.ServiceManager{
		[]*service.Service{},
	}

	for i := 1; i <= 10; i++ {
		log.Println("Initialising service: ", i)
		service := service.New(fmt.Sprintf("service-%d", i), client)
		go service.Init()
		http.HandleFunc(fmt.Sprintf("/%d/hello", i), service.Handler)
		serviceManager.Add(service)
	}

	http.HandleFunc("/list", func(writer http.ResponseWriter, request *http.Request) {
		res := fmt.Sprintf("Total services: %d\n", len(serviceManager.Services))
		for _, s := range serviceManager.Services {
			res = res + fmt.Sprintf("%s - %s\n", s.ID, s.Status)
		}
		writer.Write([]byte(res))
	})

	log.Println("Starting HTTP server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
