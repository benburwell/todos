package main

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
