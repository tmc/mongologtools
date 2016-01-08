package main

import (
	"errors"
	"io"
	"net/url"
)

type IO interface {
	Reader() (io.Reader, error)
	Writer() (io.Writer, error)
}

type InitIO func(path string) IO

type registry map[string]InitIO

var (
	ErrAlreadyRegistered = errors.New("io: already registered")
	ErrNotRegistered     = errors.New("io: not registered")
)

var r registry

func init() {
	r = registry(map[string]InitIO{})
}

func Register(scheme string, initFn InitIO) error {
	if _, ok := r[scheme]; ok {
		return ErrAlreadyRegistered
	}
	r[scheme] = initFn
	return nil
}

func Get(source string) (IO, error) {
	path, err := url.Parse(source)
	if err != nil {
		return nil, err
	}
	sourceName, sourcePath := path.Scheme, path.Host+path.Path

	fn, ok := r[sourceName]
	if !ok {
		return nil, ErrNotRegistered
	}
	return fn(sourcePath), nil
}
