package url

import (
	"github.com/paoloposso/url_shrt/util"
)

type IService interface {
	ShortenURL(string) (string, error)
	GetUrl(string) (string, error)
}

type Service struct {
	repository    IRepository
	configService util.ConfigService
}

// NewService creates a new instance of Service
func NewService(repository IRepository, config util.ConfigService) IService {
	return Service{
		repository:    repository,
		configService: config,
	}
}

func (s Service) ShortenURL(url string) (string, error) {
	shortened := encode(url)

	if err := s.repository.Save(shortened, url); err != nil {
		return "", err
	}

	baseURL := s.configService.GetBaseURL()

	return baseURL + shortened, nil
}

func (s Service) GetUrl(shortURL string) (string, error) {
	return s.repository.Find(shortURL)
}
