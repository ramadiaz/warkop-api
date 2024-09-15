package services

import "warkop-api/dto"

func (s *compServices) RegisterTransaction(data dto.Transaction) (*dto.Transaction, error) {
	id, err := s.repo.RegisterTransaction(data)
	if err != nil {
		return nil, err
	}

	for _, item := range data.Menus {
		item.TransactionID = *id

		err := s.repo.RegisterTransactionItem(*item)
		if err != nil {
			return nil, err
		}
	}

	data.ID = *id

	return &data, nil
}

func (s *compServices) GetTransaction(id string) (*dto.Transaction, error) {
	data, err := s.repo.GetTransaction(id)
	if err != nil {
		return nil, err
	}

	item_data, err := s.repo.GetTransactionItem(id)
	if err != nil {
		return nil, err
	}

	data.Menus = item_data
	
	return data, nil
}
