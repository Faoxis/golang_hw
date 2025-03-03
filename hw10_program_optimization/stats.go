package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
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
	users, errors := getUsers(r)
	result := make(DomainStat)
	for {
		select {
		case err := <-errors:
			return nil, fmt.Errorf("get users error: %w", err)
		case user, ok := <-users:
			if !ok {
				return result, nil
			}
			if strings.HasSuffix(strings.ToLower(user.Email), domain) {
				result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
			}
		}
	}
}

func getUsers(r io.Reader) (<-chan User, <-chan error) {
	scanner := bufio.NewScanner(r)
	users := make(chan User, 100)
	errors := make(chan error)

	wg := sync.WaitGroup{}
	for scanner.Scan() {
		wg.Add(1)
		bytes := scanner.Bytes()
		go func(bytes []byte) {
			defer wg.Done()
			var user User
			if err := json.Unmarshal(bytes, &user); err != nil {
				errors <- fmt.Errorf("unmarshal error: %w", err)
			}
			users <- user
		}(append([]byte{}, bytes...))
	}
	go func() {
		defer close(users)
		defer close(errors)
		wg.Wait()
	}()
	return users, errors
}
