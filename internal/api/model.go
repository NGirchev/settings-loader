package api

import (
	"bytes"
)

type Request struct {
	Type    string // default core
	Version string // default 1.0.0
	Hash    []byte // expected hash
}

type Response struct {
	Type    string // same as in request
	Version string // same as in request
	Hash    []byte // new hash
	Content []byte // content is nil if 'expected hash' != 'new hash'
}

func (r1 Response) Equal(r2 Response) bool {
	return Equals(r1, r2)
}

func Equals(r1 Response, r2 Response) bool {
	return r1.Type == r2.Type &&
		r1.Version == r2.Version &&
		bytes.Equal(r1.Hash, r2.Hash) &&
		bytes.Equal(r1.Content, r2.Content)
}
