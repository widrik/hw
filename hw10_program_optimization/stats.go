package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
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
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

type users []string

var (
	user User
	validUsersCount int
)

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)

	validUsersCount = 0
	for scanner.Scan() {
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return
		}
		result = append(result, user.Email)
		validUsersCount++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domainRegexp := regexp.MustCompile("\\."+domain)

	for _, userEmail := range u[:validUsersCount] {
		if userEmail != "" {
			if domainRegexp.MatchString(userEmail) {
				result[strings.ToLower(strings.SplitN(userEmail, "@", 2)[1])]++
			}
		}
	}
	return result, nil
}
