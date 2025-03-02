package gateways

type FakeCustomersGateway struct {
	CustomersDTO []CustomerDTO
}

func (f *FakeCustomersGateway) Create(customerDTO CustomerDTO) error {
	f.CustomersDTO = append(f.CustomersDTO, customerDTO)
	return nil
}

func (f *FakeCustomersGateway) FindOneByEmail(email string) (*CustomerDTO, error) {
	for _, customerDTO := range f.CustomersDTO {
		if customerDTO.Email == email {
			return &customerDTO, nil
		}
	}

	return nil, nil
}

func (f *FakeCustomersGateway) ExistsByEmail(email string) (bool, error) {
	for _, customerDTO := range f.CustomersDTO {
		if customerDTO.Email == email {
			return true, nil
		}
	}

	return false, nil
}
