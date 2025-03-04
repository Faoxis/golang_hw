package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"

	//nolint:depguard
	"github.com/buger/jsonparser"
)

const EmailKey = "Email"

type DomainStat map[string]int

func merge(stats []DomainStat) DomainStat {
	result := make(DomainStat)
	for _, stat := range stats {
		for k, v := range stat {
			result[k] += v
		}
	}
	return result
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	emails, errors := getUsers(r)

	wg := sync.WaitGroup{}
	workerCount := runtime.NumCPU() * 2
	stats := make([]DomainStat, workerCount)
	for i := 0; i < workerCount; i++ {
		stats[i] = make(DomainStat)
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			localStats := stats[workerId]
			for email := range emails {
				if strings.HasSuffix(strings.ToLower(email), "."+domain) {
					domainName := strings.ToLower(strings.SplitN(email, "@", 2)[1])
					localStats[domainName]++
				}
			}
		}(i)
	}
	wg.Wait()

	for err := range errors {
		return nil, fmt.Errorf("get users error: %w", err)
	}

	return merge(stats), nil
}

func getUsers(r io.Reader) (<-chan string, <-chan error) {
	scanner := bufio.NewScanner(r)

	var wg sync.WaitGroup

	workerCount := runtime.NumCPU() * 4
	jobs := make(chan []byte, workerCount)
	go func() {
		defer close(jobs)
		for scanner.Scan() {
			jobs <- append([]byte(nil), scanner.Bytes()...)
		}
	}()

	emails := make(chan string, workerCount)
	errors := make(chan error, 1)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for data := range jobs {
				email, err := jsonparser.GetString(data, EmailKey)
				if err != nil {
					errors <- fmt.Errorf("unmarshal error: %w", err)
					return
				}
				emails <- email
			}
		}()
	}
	go func() {
		wg.Wait()
		close(emails)
		close(errors)
	}()

	return emails, errors
}
