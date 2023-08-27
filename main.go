package main

import (
	"github.com/wadeling/binary-parser/pkg"
	"log"
	"os"
	"regexp"
)

var redisRule = `(?s)payload %5.*(?P<version>\d.\d\.\d\d*?)[a-z0-9]{12}-[0-9]{19}`
var redisReg = regexp.MustCompile(redisRule)

func main() {
	log.Println("start")

	testFile := "./data/redis-sentinel"
	data, err := os.ReadFile(testFile)
	if err != nil {
		log.Fatalf("failed to read file.%v", err)
		return
	}

	log.Printf("data: %v", len(data))

	log.Printf("subexp name:%v,len %v", redisReg.SubexpNames(), len(redisReg.SubexpNames()))

	// parse
	allMatches := MatchNamedCaptureGroups(redisReg, string(data))
	log.Printf("all matches:%v", allMatches)

	// check match
	_ = pkg.MatchVuln()

	log.Printf("end")
}

func MatchNamedCaptureGroups(regEx *regexp.Regexp, content string) map[string]string {
	// note: we are looking across all matches and stopping on the first non-empty match. Why? Take the following example:
	// input: "cool something to match against" pattern: `((?P<name>match) (?P<version>against))?`. Since the pattern is
	// encapsulated in an optional capture group, there will be results for each character, but the results will match
	// on nothing. The only "true" match will be at the end ("match against").
	allMatches := regEx.FindAllStringSubmatch(content, -1)
	var results map[string]string
	for _, match := range allMatches {

		// fill a candidate results map with named capture group results, accepting empty values, but not groups with
		// no names
		for nameIdx, name := range regEx.SubexpNames() {
			if nameIdx > len(match) || len(name) == 0 {
				continue
			}
			if results == nil {
				results = make(map[string]string)
			}
			results[name] = match[nameIdx]
		}
		// note: since we are looking for the first best potential match we should stop when we find the first one
		// with non-empty results.
		if !isEmptyMap(results) {
			break
		}
	}
	return results
}
func isEmptyMap(m map[string]string) bool {
	if len(m) == 0 {
		return true
	}
	for _, value := range m {
		if value != "" {
			return false
		}
	}
	return true
}
