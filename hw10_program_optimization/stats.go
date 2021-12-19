package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type UserIterator struct {
	scanner *bufio.Scanner
	jit     jsoniter.API
	user    *User
}

func (i *UserIterator) hasNext() bool {
	return i.scanner.Scan()
}

func (i *UserIterator) next() (*User, error) {
	err := i.jit.Unmarshal(i.scanner.Bytes(), i.user)
	if err != nil {
		return nil, err
	}

	i.user.Email = strings.ToLower(i.user.Email)
	return i.user, err
}

func NewIterator(reader io.Reader) *UserIterator {
	return &UserIterator{
		scanner: bufio.NewScanner(reader),
		jit:     jsoniter.ConfigCompatibleWithStandardLibrary,
		user:    &User{},
	}
}

func getUsers(r io.Reader) (*UserIterator, error) { //nolint:unparam
	return NewIterator(r), nil
}

func countDomains(uit *UserIterator, domain string) (DomainStat, error) {
	result := make(DomainStat)

	find := "." + domain

	var i int
	for uit.hasNext() {
		user, err := uit.next()
		if err != nil {
			return nil, err
		}

		if i = strings.LastIndex(user.Email, "@"); i == -1 {
			return nil, fmt.Errorf("string is not correct email - %s", user.Email)
		}

		if !strings.HasSuffix(user.Email, find) {
			continue
		}

		result[user.Email[i+1:]]++
	}

	return result, nil
}
