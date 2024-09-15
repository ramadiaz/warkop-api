package services

import "warkop-api/dto"

func (s *compServices) RegisterMenu(data dto.Menu) error {
	return s.repo.RegisterMenu(data)
}

func (s *compServices) GetAllMenu() ([]*dto.Menu, error) {
	return s.repo.GetAllMenu()
}
