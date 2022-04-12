package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNotSavedPages = fmt.Errorf("no saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	_, err := io.WriteString(h, p.URL)
	if err != nil {
		return "", fmt.Errorf("failed to calculate URL hash due error: %v", err)
	}

	_, err = io.WriteString(h, p.UserName)
	if err != nil {
		return "", fmt.Errorf("failed to calculate UserName hash due error: %v", err)
	}

	return string(h.Sum(nil)), nil
}
