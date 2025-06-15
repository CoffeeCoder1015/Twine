package main

import "regexp"

type queryFilterPattern struct {
	directory string
	name      *regexp.Regexp
	fileSize  string
	mode      *regexp.Regexp
	date      string
	DirFile   string
}
