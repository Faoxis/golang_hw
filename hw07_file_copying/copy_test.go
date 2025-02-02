package main

import (
	"fmt"
	"testing"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("Test coping file with zero offset", func(t *testing.T) {
		err := Copy(
			"./testdata/input.txt",
			"/tmp/result.txt",
			0,
			1000000,
		)
		fmt.Println(err)
	})
}
