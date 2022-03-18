package reaper

import "doublequote/pkg/domain"

type Service struct {
	entryService *domain.EntryService
}

func New(entryService *domain.EntryService) *Service {
	return &Service{entryService}
}

func (s *Service) Run() {
	// TODO reaper 💀
}
