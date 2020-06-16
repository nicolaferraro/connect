package storage

import (
	"github.com/nicolaferraro/connect/pkg/token"
)

type TokenStorage interface {
	TokenLister
	TokenGetter
	TokenSaver
}

type TokenLister interface {
	List() ([]string, error)
}

type TokenGetter interface {
	Get(name string) (*token.Token, error)
}

type TokenSaver interface {
	Save(string, *token.Token) error
}
