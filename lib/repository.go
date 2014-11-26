package ghcloc

import (
	"errors"
	"strings"
)

type Repository struct {
	User     string
	Name     string
	UserAuth bool
	AuthUser string
	AuthPass string
	ReqSem   chan struct{}
}

func NewRepository(user string, name string) *Repository {
	return &Repository{user, name, false, "", "", make(chan struct{}, 100)}
}

func ParseRepository(compacted string) (*Repository, error) {
	comps := strings.Split(compacted, "/")
	if len(comps) != 2 {
		err := errors.New("Repository name must be of the form 'user/repo'")
		return nil, err
	}
	return NewRepository(comps[0], comps[1]), nil
}

func (self *Repository) String() string {
	return self.User + "/" + self.Name
}

func (self *Repository) Authenticate(username string, password string) {
	self.UserAuth = true
	self.AuthUser = username
	self.AuthPass = password
}
