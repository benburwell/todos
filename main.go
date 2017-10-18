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

func GetErrorMessages(errs []error) []string {
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = err.Error()
	}
	return messages
}

func main() {
	nameArg := flag.String("name", "", "Name to look for")
	helpArg := flag.Bool("help", false, "Get help")
	flag.Parse()

	if *helpArg {
		flag.Usage()
		os.Exit(0)
	}

	config := GetConfig()
	isValid, errors := config.Validate()

	if fmt.Sprintf("%s", *nameArg) != "" {
		config.Name = *nameArg
	}

	if !isValid {
		die(strings.Join(GetErrorMessages(errors), "\n"))
	}

	todofiles, err := scanDir(".", config)
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

func scanTodo(reader *bufio.Reader, config *Config) bool {
	patternScanner, err := NewPatternScanner(fmt.Sprintf("TODO(%s): ", config.Name))
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

func scanFile(filepath string, config *Config) ([]Todo, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	var todos []Todo
	for scanTodo(reader, config) {
		todo, err := readTodo(reader)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func scanDir(dir string, config *Config) ([]TodoFile, error) {
	fileinfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []TodoFile
	for _, fileinfo := range fileinfos {
		if fileinfo.Mode().IsRegular() && config.ShouldScanFile(fileinfo.Name()) {
			filepath := path.Join(dir, fileinfo.Name())
			todos, err := scanFile(filepath, config)
			if err != nil {
				return nil, err
			}
			if len(todos) > 0 {
				todofile := TodoFile{}
				todofile.File = filepath
				todofile.Todos = todos
				result = append(result, todofile)
			}
		} else if fileinfo.IsDir() && config.ShouldScanDir(fileinfo.Name()) {
			dirname := path.Join(dir, fileinfo.Name())
			infos, err := scanDir(dirname, config)
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
