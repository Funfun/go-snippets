package main

import (
	"database/sql"
	"sync"
)

type User struct{ active bool }

func (u User) Active() bool { return u.active }

type Users interface {
	FindAll() []*User
}
type API struct {
	Users
}

func (api *API) ActiveUsers() []*User {
	users := api.FindAll()
	active := []*User{}

	for _, user := range users {
		if user.Active() {
			active = append(active, user)
		}
	}

	return active
}

type KVStorage struct {
	sync.RWMutex
	storage map[interface{}]interface{}
}

func (kv *KVStorage) FindAll() []*User {
	kv.RLock()
	defer kv.RUnlock()

	users := []*User{}
	for _, v := range kv.storage {
		user, ok := v.(*User)
		if ok {
			users = append(users, user)
		}
	}

	return users
}

type MySQLStorage struct {
	db *sql.DB
}

func (s *MySQLStorage) FindAll() []*User {
	users := []*User{}
	// use s.db here

	return users
}
