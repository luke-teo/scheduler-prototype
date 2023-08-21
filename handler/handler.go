package handler

import (
	"github.com/scheduler-prototype/mgraph"
	"github.com/scheduler-prototype/repository"
)

type Handler struct {
	client *mgraph.MGraph
	repo   *repository.Repository
}

func NewHandler(client *mgraph.MGraph, repo *repository.Repository) *Handler {
	return &Handler{
		client: client,
		repo:   repo,
	}
}
