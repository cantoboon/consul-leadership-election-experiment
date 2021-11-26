package service

type ServiceManager struct {
	Services []*Service
}

func (sm *ServiceManager) Add(service *Service) {
	sm.Services = append(sm.Services, service)
}
