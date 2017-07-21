package types

import (
	"io"
	"net/http"
)

type HookContext interface {
	Delete(string) error
	GetByHook(string) (*HookDal, error)
	GetByTypeUrl(string, string) (*HookDal, error)
	GetAll() ([]string, []*HookDal, error)
	Errorf(string, ...interface{})
	Put(*HookDal) error
	UrlPost(*HookDal, io.Reader) (*http.Response, error)
}
