package url_test

import (
	"testing"
	"time"

	"github.com/paoloposso/url_shrt/url"
)

// MockRepository is a mock implementation of IRepository
type MockRepository struct{}

func (mr *MockRepository) Save(shortURL string, longURL string) error {
	return nil
}

func (mr *MockRepository) Find(shortURL string) (string, error) {
	return "http://www.google.com", nil
}

// MockConfigService is a mock implementation of util.ConfigService
type MockConfigService struct{}

func (mcs *MockConfigService) GetBaseURL() string {
	return "http://mockbaseurl.com/"
}

func (mcs MockConfigService) GetMongoDbTimeOut() time.Duration {
	return 10 * time.Second
}

func TestShortenURL(t *testing.T) {
	repo := &MockRepository{}
	configService := &MockConfigService{}
	service := url.NewService(repo, configService)

	result, err := service.ShortenURL("http://www.google.com")

	if err != nil || result == "" {
		t.Error("Error while shortening URL")
	}
}

func TestGetUrl(t *testing.T) {
	repo := &MockRepository{}
	configService := &MockConfigService{}
	service := url.NewService(repo, configService)

	result, err := service.GetUrl("http://mockbaseurl.com/123456")

	if err != nil || result == "" {
		t.Error("Error while getting URL")
	}
}
