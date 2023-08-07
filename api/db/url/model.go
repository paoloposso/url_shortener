package db_url

type URLDocument struct {
	ShortURL string `bson:"shortUrl"`
	URL      string `bson:"url"`
}
