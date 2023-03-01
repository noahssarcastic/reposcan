package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	args := os.Args[1:]
	rootDir := args[0]

	re, err := regexp.Compile(`ENV\[['"](\w+)['"]\]`)
	check(err)

	parseDir(rootDir, func(line string) []string {
		return findAll(line, re)
	})
}

func findAll(s string, re *regexp.Regexp) (matches []string) {
	results := re.FindAllStringSubmatch(s, -1)
	if results == nil {
		return
	}
	for _, capture := range results {
		matches = append(matches, capture[1])
	}
	return matches
}

type Matcher = func(string) []string

func parseDir(name string, matcher Matcher) {
	files, err := os.ReadDir(name)
	check(err)
	path, err := filepath.Abs(name)
	check(err)
	for _, child := range files {
		childName := filepath.Join(path, child.Name())
		if child.IsDir() {
			parseDir(childName, matcher)
			continue
		}
		handleFile(childName, matcher)
	}
}

// type Report = struct {
// 	fileName string
// 	matches  []Match
// }

// type Match = struct {
// 	s          string
// 	lineNumber int
// }

func handleFile(name string, matcher Matcher) {
	file, err := os.Open(name)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// For each line...
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		matches := matcher(line)
		if len(matches) > 0 {
			fmt.Printf("%s:%d\n%v\n", name, i, printMatches(matches))
		}
	}
	check(scanner.Err())
}

func printMatches(matches []string) string {
	return strings.Join(matches, "\n")
}
