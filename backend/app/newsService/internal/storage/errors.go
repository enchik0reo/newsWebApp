package storage

import "errors"

var (
	ErrNoSources        = errors.New("there are no sources")
	ErrSourceNotFound   = errors.New("source not found")
	ErrSourceExists     = errors.New("source exists")
	ErrNoNewArticles    = errors.New("there are no new articles")
	ErrNoLatestArticles = errors.New("there are no latest articles")
)