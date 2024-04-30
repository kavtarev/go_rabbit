package main

type Storage interface {
	GetSome() string
}

type MockStorage struct {
	
}

func (s *MockStorage) GetSome() string {
	return "get-some"
}