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

var includeExtensions = map[string]bool{
	"css":  true,
	"go":   true,
	"java": true,
	"js":   true,
	"scss": true,
	"ts":   true,
}

var ignoreDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
}

type TodoFile struct {
	File  string
	Todos []Todo
}

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

func die(message string) {
	os.Stderr.WriteString(strings.Join([]string{message, "\n"}, ""))
	os.Exit(1)
}

func shouldScanDir(dir string) bool {
	_, ok := ignoreDirs[dir]
	return ok == false
}

func shouldScanFile(f os.FileInfo) bool {
	parts := strings.Split(f.Name(), ".")
	if len(parts) < 2 {
		return false
	}
	ext := parts[len(parts)-1]
	_, ok := includeExtensions[ext]
	return ok
}

func makePattern(pattern string) (byte, string) {
	return byte(pattern[0]), pattern[1:]
}

func scanTodo(reader *bufio.Reader, name string) bool {
	startByte, restString := makePattern(fmt.Sprintf("TODO(%s): ", name))
	_, err := reader.ReadBytes(startByte)
	if err != nil {
		return false
	}
	next, err := reader.Peek(10)
	if string(next) == restString {
		reader.Discard(len(restString))
		return true
	}
	return scanTodo(reader, name)
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
