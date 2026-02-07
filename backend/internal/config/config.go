package config

import "os"

type Config struct {
	TheMovieDBToken string `env:"THE_MOVIE_DB_TOKEN"`
}

func LoadConfig() (config Config, err error) {
	cfg := Config{
		TheMovieDBToken: os.Getenv("THE_MOVIE_DB_TOKEN"),
	}
	return cfg, nil
}
