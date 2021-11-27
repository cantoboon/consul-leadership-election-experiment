// Package service
//
// Defines a "service". A service has an HTTP handler, so we can get some information out of it.
// On calling Init(), it will create a session and attempt to get the lock. Status of which can be
// viewed with the HTTP endpoint.
package service

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"time"
)

// Status of the Service.
type Status string

const (
	UNINITIALISED Status = "Uninitialized"
	LEADER        Status = "Leader"
	FOLLOWER      Status = "Follower"
)

const (
	KEY = "service/leader"
)

type Service struct {
	ID              string
	consulSessionId string
	consulClient    *api.Client
	Status          Status
	Leader          string
}

func New(id string, client *api.Client) *Service {
	return &Service{
		ID:              id,
		consulSessionId: "",
		consulClient:    client,
		Status:          UNINITIALISED,
	}
}

func (s *Service) Init() error {
	session, err := s.startSession()
	if err != nil {
		return err
	}
	s.consulSessionId = session

	err = s.registerService()
	if err != nil {
		return err
	}

	err = s.acquireLeadership()
	if err != nil {
		return err
	}

	log.Println("Initialised service: ", s.ID)
	return nil
}

func (s *Service) startSession() (string, error) {
	log.Println("Starting session: ", s.ID)
	sessionClient := s.consulClient.Session()

	sessionId, _, err := sessionClient.Create(
		&api.SessionEntry{
			Name: fmt.Sprintf(s.ID),
			TTL:  "10s",
		},
		nil,
	)
	if err != nil {
		return "", err
	}

	go func() {
		for {
			_, _, err := sessionClient.Renew(sessionId, nil)
			if err != nil {
				log.Println("Failed to renew session for service: ", s.ID, err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	log.Println("Session started: ", s.ID, "Session ID: ", sessionId)
	return sessionId, nil
}

func (s *Service) registerService() error {
	return nil
}

// Attempts to become the leader
func (s *Service) acquireLeadership() error {
	kvClient := s.consulClient.KV()

	kvPair := &api.KVPair{Session: s.consulSessionId, Key: KEY, Value: []byte(s.ID)}

	res, _, err := kvClient.Acquire(kvPair, nil)
	if err != nil {
		return err
	}

	if res {
		log.Println(s.ID, " is the leader")
		s.Status = LEADER
		s.Leader = s.ID
	} else {
		s.Status = FOLLOWER
		leader, _, err := kvClient.Get(KEY, nil)
		if err != nil {
			log.Println("Failed to get current leader")
		}
		s.Leader = string(leader.Value)
	}
	return nil
}

func (s *Service) Handler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(fmt.Sprintf("Hello I am: %s\nStatus: %s\nCurrent Leader: %s", s.ID, s.Status, s.Leader)))
}
