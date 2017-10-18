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
		Die(strings.Join(GetErrorMessages(errors), "\n"))
	}

	patternScanner, err := NewPatternScanner(fmt.Sprintf("TODO(%s): ", config.Name))
	if err != nil {
		Die(fmt.Sprintf("Error creating scanner: %s", err.Error()))
	}
	todofiles, err := scanDir(".", config, patternScanner)
	if err != nil {
		Die(fmt.Sprintf("Error reading files: %s", err.Error()))
	}

	for _, f := range todofiles {
		fmt.Println(f.File)
		for _, t := range f.Todos {
			fmt.Println("  - ", t.Summary)
		}
		fmt.Println()
	}
}

func scanFile(filepath string, scanner *PatternScanner) ([]Todo, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	var todos []Todo
	for scanner.Scan(reader) {
		todo, err := scanner.Read(reader)
		if err != nil {
			return nil, err
		}
		todos = append(todos, *todo)
	}
	return todos, nil
}

func scanDir(dir string, config *Config, scanner *PatternScanner) ([]TodoFile, error) {
	fileinfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []TodoFile
	for _, fileinfo := range fileinfos {
		if fileinfo.Mode().IsRegular() && config.ShouldScanFile(fileinfo.Name()) {
			filepath := path.Join(dir, fileinfo.Name())
			todos, err := scanFile(filepath, scanner)
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
			infos, err := scanDir(dirname, config, scanner)
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
