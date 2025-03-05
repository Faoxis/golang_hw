package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	//nolint:depguard
	"github.com/buger/jsonparser"
)

const EmailKey = "Email"

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	stat := DomainStat{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		email, err := jsonparser.GetString(line, EmailKey)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(strings.ToLower(email), "."+domain) {
			domainName := strings.ToLower(strings.SplitN(email, "@", 2)[1])
			stat[domainName]++
		}
	}
	return stat, nil
}
