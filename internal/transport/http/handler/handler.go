package handler

import (
	"github.com/6ivkin/test.git/internal/reader"
)

type ReaderHandler struct {
	service *reader.Service
}

func NewReaderHandler(
	service *reader.Service,
) *ReaderHandler {
	return &ReaderHandler{
		service: service,
	}
}
