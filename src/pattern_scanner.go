package main

import (
	"bufio"
	"errors"
	"strings"
)

type PatternScanner struct {
	firstByte byte
	rest      string
}

func NewPatternScanner(pattern string) (*PatternScanner, error) {
	if len(pattern) < 2 {
		return nil, errors.New("Pattern must be at least two characters")
	}
	fb := byte(pattern[0])
	r := pattern[1:]
	return &PatternScanner{firstByte: fb, rest: r}, nil
}

func (p *PatternScanner) Scan(reader *bufio.Reader) bool {
	_, err := reader.ReadBytes(p.firstByte)
	if err != nil {
		return false
	}
	next, err := reader.Peek(len(p.rest))
	if err != nil {
		return false
	}
	if string(next) == p.rest {
		reader.Discard(len(p.rest))
		return true
	}
	return p.Scan(reader)
}

func (p *PatternScanner) Read(reader *bufio.Reader) (*Todo, error) {
	bytes, err := reader.ReadBytes(byte('\n'))
	if err != nil {
		return nil, err
	}
	todo := &Todo{
		Summary: strings.TrimSpace(string(bytes)),
	}
	return todo, nil
}
