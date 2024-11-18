package types

import (
	gh "ghclient/ghclient"
	"sync"
)

type UserStorage struct {
	mu    sync.Mutex
	users []gh.UserFormattedData
}

func (us *UserStorage) Append(user gh.UserFormattedData) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.users = append(us.users, user)
}

func (us *UserStorage) Value() []gh.UserFormattedData {
	us.mu.Lock()
	defer us.mu.Unlock()
	return us.users
}
