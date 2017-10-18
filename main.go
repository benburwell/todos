package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Only consider reading files that end with one of these extensions preceeded
// by a period (e.g. `.go`)
var includeExtensions = map[string]bool{
	"css":  true,
	"go":   true,
	"java": true,
	"js":   true,
	"scss": true,
	"ts":   true,
}

// Don't look in any directories with these names
var ignoreDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
}

// A TodoFile consists of a path relative to the running directory and a
// collection of Todos that are contained in it
type TodoFile struct {
	File  string
	Todos []Todo
}

// A Todo holds a Summary
//
// In the future, it may also hold other information
type Todo struct {
	Summary string
}

func main() {
	nameArg := flag.String("name", "", "Name to look for")
	helpArg := flag.Bool("help", false, "Get help")
	flag.Parse()

	if *helpArg {
		flag.Usage()
		os.Exit(0)
	}

	if fmt.Sprintf("%s", *nameArg) == "" {
		die("Name must be specified")
	}

	todofiles, err := scanDir(".", *nameArg)
	if err != nil {
		die(fmt.Sprintf("Error reading files: %s", err.Error()))
	}

	for _, f := range todofiles {
		fmt.Println(f.File)
		for _, t := range f.Todos {
			fmt.Println("  - ", t.Summary)
		}
	}
}

// Print a message out to standard error and exit the program
func die(message string) {
	os.Stderr.WriteString(strings.Join([]string{message, "\n"}, ""))
	os.Exit(1)
}

// Determine whether the specified directory name should be scanned
func shouldScanDir(dir string) bool {
	_, ok := ignoreDirs[dir]
	return ok == false
}

// Determine whether the specified file should be scanned
func shouldScanFile(f os.FileInfo) bool {
	parts := strings.Split(f.Name(), ".")
	if len(parts) < 2 {
		return false
	}
	ext := parts[len(parts)-1]
	_, ok := includeExtensions[ext]
	return ok
}

func scanTodo(reader *bufio.Reader, name string) bool {
	patternScanner, err := NewPatternScanner(fmt.Sprintf("TODO(%s): ", name))
	if err != nil {
		return false
	}
	return patternScanner.Scan(reader)
}

func readTodo(reader *bufio.Reader) (Todo, error) {
	bytes, err := reader.ReadBytes(byte('\n'))
	if err != nil {
		return Todo{}, err
	}
	return Todo{Summary: string(bytes)}, nil
}

func scanFile(filepath string, name string) ([]Todo, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	var todos []Todo
	for scanTodo(reader, name) {
		todo, err := readTodo(reader)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func scanDir(dir string, name string) ([]TodoFile, error) {
	fileinfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []TodoFile
	for _, fileinfo := range fileinfos {
		if fileinfo.Mode().IsRegular() && shouldScanFile(fileinfo) {
			filepath := path.Join(dir, fileinfo.Name())
			todos, err := scanFile(filepath, name)
			if err != nil {
				return nil, err
			}
			if len(todos) > 0 {
				todofile := TodoFile{}
				todofile.File = filepath
				todofile.Todos = todos
				result = append(result, todofile)
			}
		} else if fileinfo.IsDir() && shouldScanDir(fileinfo.Name()) {
			dirname := path.Join(dir, fileinfo.Name())
			infos, err := scanDir(dirname, name)
			if err != nil {
				return nil, err
			}
			for _, info := range infos {
				result = append(result, info)
			}
		}
	}
	return result, nil
}
