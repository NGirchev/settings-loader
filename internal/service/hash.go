package service

import (
	"crypto/md5"
	"crypto/sha256"
)

type IHasher interface {
	Hash(data []byte) []byte
}

func NewMD5Hasher() IHasher {
	return &MD5Hasher{}
}

func NewSHA256Hasher() IHasher {
	return &SHA256Hasher{}
}

type MD5Hasher struct{}
type SHA256Hasher struct{}

func (h *MD5Hasher) Hash(data []byte) []byte {
	hasher := md5.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func (h *SHA256Hasher) Hash(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
