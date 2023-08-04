package url

import (
	"os"
)

type IService interface {
	ShortenURL(string) (string, error)
	GetUrl(string) (string, error)
}

type Service struct {
	repository IRepository
}

// NewService creates a new instance of Service
func NewService(repository IRepository) IService {
	return Service{
		repository: repository,
	}
}

func (s Service) ShortenURL(url string) (string, error) {
	shortened := encode(url)

	if err := s.repository.Save(shortened, url); err != nil {
		return "", err
	}

	baseURL := os.Getenv("BASE_URL")

	return baseURL + shortened, nil
}

func (s Service) GetUrl(shortURL string) (string, error) {
	return s.repository.Find(shortURL)
}
