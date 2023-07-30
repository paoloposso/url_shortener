package url

import (
	"crypto/sha256"
	"math/big"
	"os"
)

type IService interface {
	ShortenURL(string) (string, error)
}

type Service struct {
	repository IRepository
}

// NewService is a function that returns a pointer to a Service struct.
func NewService(repository IRepository) Service {
	return Service{
		repository: repository,
	}
}

// ShortenURL is a method on the Service struct.
// It takes a string and returns a string.
func (s Service) ShortenURL(url string) (string, error) {
	shortened := s.shortenURL(url)

	if err := s.repository.Save(url, shortened); err != nil {
		return "", err
	}

	baseURL := os.Getenv("BASE_URL")

	return baseURL + shortened, nil
}

func (s Service) shortenURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	number := new(big.Int).SetBytes(hash[:])

	return s.base62Encode(number)
}

func (s Service) base62Encode(number *big.Int) string {
	const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	base := big.NewInt(int64(len(base62Chars)))
	zero := big.NewInt(0)
	encoded := ""

	for number.Cmp(zero) > 0 {
		mod := new(big.Int)
		number.DivMod(number, base, mod)

		encoded = string(base62Chars[mod.Int64()]) + encoded
	}

	return encoded
}
