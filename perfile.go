package main

import (
	"bufio"
	"os"
	"path/filepath"
)

type PerFile struct {
	name    string
	file    *os.File
	scanner *bufio.Scanner
	cyclic  bool
	eof     bool
}

func NewPerFile(name string, cyclic bool) (*PerFile, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return &PerFile{
		name:    filepath.Base(name),
		file:    file,
		scanner: bufio.NewScanner(file),
		cyclic:  cyclic,
		eof:     false,
	}, nil
}

func (p *PerFile) Next() bool {
	ret := p.scanner.Scan()
	if !ret {
		if p.cyclic {
			_, err := p.file.Seek(0, 0)
			if err != nil {
				p.eof = true
				return false
			}
			p.scanner = bufio.NewScanner(p.file)
			if !p.scanner.Scan() {
				p.eof = true
			}
		} else {
			p.eof = true
		}
	}

	return ret
}

func (p *PerFile) Get() string {
	if p.eof {
		return ""
	}

	return p.scanner.Text()
}
