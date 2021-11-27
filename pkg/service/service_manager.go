package service

type Manager struct {
	Services []*Service
}

func (sm *Manager) Add(service *Service) {
	sm.Services = append(sm.Services, service)
}
