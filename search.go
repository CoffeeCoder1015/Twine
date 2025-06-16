package main

import (
	"io/fs"
	"path/filepath"
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

type resultEntry struct{
    fs.DirEntry
    path string 
}

type Twine struct{
    filter queryFilterPattern
    cache map[string][]resultEntry
}

func InitTwine() Twine{
    t := Twine{
        filter: queryFilterPattern{
            directory:".",
        },
        cache: make(map[string][]resultEntry),
    }
    return t
}

func (t Twine) Query() []resultEntry{
    t.Search()
    return t.cache[t.filter.directory]
}

func (t Twine) Search(){
    _,c := t.cache[t.filter.directory]
    if c {
        return
    }
    
    results := make([]resultEntry,0)
    filepath.WalkDir(t.filter.directory,func(epath string, d fs.DirEntry, err error) error {
        if err == nil {
            results = append(results, resultEntry{
                path: filepath.Dir(epath),
                DirEntry: d,
            })
        }
        return nil
    })
    t.cache[t.filter.directory] = results
}

