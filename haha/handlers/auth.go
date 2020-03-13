package handlers

import (
	"joblessness/haha/models"
	"sync"
)

type AuthHandler struct {
	Sessions    map[string]uint64
	Users       map[string]*models.Person
	UserAvatars map[uint64]string
	UserSummary map[uint64]models.UserSummary
	Mu          sync.RWMutex
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler {
		Sessions: make(map[string]uint64, 10),
		Users:    map[string]*models.Person{
			"marat1k": {1, "marat1k", "ABCDE12345", "Marat", "Ishimbaev", "m@m.m", "89032909821"},
		},
		UserAvatars: map[uint64]string{},
		UserSummary: map[uint64]models.UserSummary{},
		Mu:          sync.RWMutex{},
	}
}