package url

type IRepository interface {
	Find(shortURL string) (string, error)
	Save(shortURL string, longURL string) error
}
