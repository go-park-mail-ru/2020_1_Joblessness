package handlers

import (
	"joblessness/haha/models"
	"sync"
)

type AuthHandler struct {
	Sessions    map[string]uint
	Users       map[string]*models.Person
	UserAvatars map[uint]string
	UserSummary map[uint]models.UserSummary
	Mu          sync.RWMutex
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler {
		Sessions: make(map[string]uint, 10),
		Users:    map[string]*models.Person{
			"marat1k": {1, "marat1k", "ABCDE12345", "Marat", "Ishimbaev", "m@m.m", "89032909821"},
		},
		UserAvatars: map[uint]string{},
		UserSummary: map[uint]models.UserSummary{},
		Mu:          sync.RWMutex{},
	}
}