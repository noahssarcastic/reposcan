package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Matcher = struct {
	glob    string
	pattern *regexp.Regexp
}

func NewMatcher(glob string, pattern *regexp.Regexp) *Matcher {
	return &Matcher{glob, pattern}
}

func main() {
	rootDir := flag.String("file", ".", "")
	glob := flag.String("glob", "*", "")
	pattern := flag.String("regexp", "", "")
	flag.Parse()

	re, err := regexp.Compile(*pattern)
	if err != nil {
		log.Fatalln(err)
	}
	matcher := NewMatcher(*glob, re)

	err = parseDir(*rootDir, matcher)
	if err != nil {
		log.Fatalln(err)
	}
}

func parseError(fileName string, cause error) error {
	return errors.Wrap(
		cause,
		fmt.Sprintf("error parsing %s", fileName),
	)
}

func parseDir(name string, matcher *Matcher) error {
	files, err := os.ReadDir(name)
	if err != nil {
		return err
	}
	path, err := filepath.Abs(name)
	if err != nil {
		return err
	}
	for _, child := range files {
		name := filepath.Join(path, child.Name())
		isDir := child.IsDir()
		err := handleChild(name, isDir, matcher)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleChild(name string, isDir bool, matcher *Matcher) error {
	if isDir {
		err := parseDir(name, matcher)
		if err != nil {
			return err
		}
		return nil
	}
	err := handleFile(name, matcher)
	if err != nil {
		log.Println(parseError(name, err))
	}
	return nil
}

// type Report = struct {
// 	fileName string
// 	matches  []Match
// }

// type Match = struct {
// 	s          string
// 	lineNumber int
// }

func handleFile(name string, matcher *Matcher) error {
	baseName := filepath.Base(name)
	if match, err := filepath.Match(matcher.glob, baseName); err != nil {
		return err
	} else if !match {
		return nil
	}

	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		matches := findAll(line, matcher.pattern)
		if len(matches) > 0 {
			fmt.Printf("%s:%d\n%v\n", name, i, strings.Join(matches, "\n"))
		}
	}
	if scanner.Err() != nil {
		return err
	}
	return nil
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
