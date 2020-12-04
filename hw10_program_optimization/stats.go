package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	stat, err := countDomainsStat(r, domain)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func countDomainsStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	var user User

	domainBytes := []byte("." + domain)

	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		if bytes.Contains(lineBytes, domainBytes) {
			user = User{}
			if err := json.Unmarshal(lineBytes, &user); err != nil {
				return nil, fmt.Errorf("get users error: %w", err)
			}
			if strings.HasSuffix(user.Email, domain) {
				d := strings.ToLower(user.Email[strings.Index(user.Email, "@")+1 : len(user.Email)])
				result[d]++
			}
		}
	}
	return result, nil
}
