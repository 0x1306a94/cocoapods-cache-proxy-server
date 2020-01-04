package config

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
)

type AuthorizationConfig struct {
	sync.RWMutex
	users []authorizationUser

	admin authorizationUser
}

type authorizationUser struct {
	user     string
	password string
	admin    bool
}

func (this *AuthorizationConfig) SetupAdminUser(user, password string) {
	this.Lock()
	defer this.Unlock()
	this.admin = authorizationUser{user, password, true}
}

func (this *AuthorizationConfig) Add(user, password string) bool {
	if len(user) == 0 || len(password) == 0 {
		return false
	}
	this.Lock()
	defer this.Unlock()
	for idx, v := range this.users {
		if v.user == user {
			this.users[idx] = authorizationUser{user, password, false}
			return true
		}
	}
	this.users = append(this.users, authorizationUser{user, password, false})
	return true
}

func (this *AuthorizationConfig) Remove(user, password string) bool {
	if len(user) == 0 || len(password) == 0 {
		return false
	}
	this.Lock()
	defer this.Unlock()
	for idx, v := range this.users {
		if v.user == user && v.password == password {
			this.users = append(this.users[:idx], this.users[idx+1:]...)
			return true
		}
	}
	return false
}

func (this *AuthorizationConfig) ValidationForBasicAuthorization(value string) bool {
	if len(value) < 7 {
		return false
	}
	if !strings.HasPrefix(value, "Basic ") {
		return false
	}
	value = strings.ReplaceAll(value, "Basic ", "")
	strByte, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Println(err)
		return false
	}
	pair := strings.Split(string(strByte), ":")
	if len(pair) != 2 {
		return false
	}
	if this.Validation(pair[0], pair[1]) || this.ValidationAdmin(pair[0], pair[1]) {
		return true
	}
	return false
}

func (this *AuthorizationConfig) Validation(user, password string) bool {
	if len(user) == 0 || len(password) == 0 {
		return false
	}
	this.RLock()
	defer this.RUnlock()
	for _, v := range this.users {
		if v.user == user && v.password == password {
			return true
		}
	}
	return false
}

func (this *AuthorizationConfig) ValidationAdmin(user, password string) bool {
	if len(user) == 0 || len(password) == 0 {
		return false
	}
	return this.admin.user == user && this.admin.password == password
}
