package api

import (
	"errors"
	"github.com/ngirchev/settings-loader/internal/service"
)

type LoaderController struct {
	loaderService service.ILoaderService
}

func NewLoaderController(loaderService service.ILoaderService) *LoaderController {
	return &LoaderController{
		loaderService: loaderService,
	}
}

func (l *LoaderController) LoadComponent(payload *Request, reply *Response) error {
	err := validate(payload)
	if err != nil {
		return err
	}

	content, hash, err := l.loaderService.LoadComponent(payload.Type, payload.Version, payload.Hash)
	if err != nil {
		return err
	}

	*reply = Response{
		Type:    payload.Type,
		Version: payload.Version,
		Hash:    hash,
		Content: content,
	}
	return nil
}

func validate(payload *Request) error {
	if payload == nil {
		return errors.New("payload is empty")
	}
	if payload.Type == "" {
		payload.Type = "core"
	}
	if payload.Version == "" {
		payload.Version = "1.0.0"
	}
	return nil
}
