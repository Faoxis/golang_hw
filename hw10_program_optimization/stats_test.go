//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("find in invalid json", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(`{"Email": }`), "com")
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("duplicate emails", func(t *testing.T) {
		duplicateEmails := `{"Email":"aliquid_qui_ea@Browsedrive.gov" }
							{"Email":"aliquid_qui_ea@Browsedrive.gov" }
							{"Email":"aliquid_qui_ea@Browsedrive.gov" }`
		result, err := GetDomainStat(bytes.NewBufferString(duplicateEmails), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 3}, result)
	})

	t.Run("case insensitive domain matching", func(t *testing.T) {
		caseData := `{"Email":"aliquid_qui_ea@example.gov" }
					 {"Email":"aliquid_qui_ea@Example.gov" }
					 {"Email":"aliquid_qui_ea@EXAMPLE.gov" }`
		result, err := GetDomainStat(bytes.NewBufferString(caseData), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"example.gov": 3}, result)
	})

	t.Run("emails without domain", func(t *testing.T) {
		caseData := `{"Email":"aliquid_qui_ea@example"}`
		result, err := GetDomainStat(bytes.NewBufferString(caseData), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("large input", func(t *testing.T) {
		inputSize := 10_000
		var largeData strings.Builder
		for i := 0; i < inputSize; i++ {
			largeData.WriteString(fmt.Sprintf(`{"Email":"aliquid_%d@example.com"}`+"\n", inputSize))
		}

		result, err := GetDomainStat(bytes.NewBufferString(largeData.String()), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"example.com": inputSize}, result)
	})
}
