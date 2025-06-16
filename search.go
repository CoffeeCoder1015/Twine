package main

import (
	"os"
	"regexp"
)

type queryFilterPattern struct {
	directory string
	name      *regexp.Regexp
	fileSize  string
	mode      *regexp.Regexp
	date      string
	DirFile   string
}

type Twine struct{
    filter queryFilterPattern
}

func (t Twine) Query() []os.DirEntry{
    entries, _ := os.ReadDir(t.filter.directory)
    return entries
}
