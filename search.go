package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/bubbles/list"
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
    cache := t.cache[t.filter.directory]
    len := min(3000,len(cache))
    return t.cache[t.filter.directory][:len]
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

func formatItems(entries []resultEntry) []list.Item{
    items :=  make([]list.Item,len(entries))
    for i, e := range entries{
        items[i] = formatInfo(e)
    } 
    return items
}

func formatInfo(entry resultEntry) item{
    info, _ := entry.Info()
    icon := "📄"
    if entry.IsDir(){
        icon = "📁"
    }
    title := fmt.Sprintf("%s %s %s",entry.Name(),entry.path,icon)
    desc := fmt.Sprintf("%s %s %s",formatSize(info.Size()),info.ModTime().Format("2006-01-02 15:04:05"),info.Mode())
    return item{
        title: title,
        desc: desc,
    }
}

func formatSize(size int64) string{
    d := int64(1)
    index := 0
    for i := 1; i < 4; i++{
        temp := d*1000
        if size > temp {
            d = temp
            index = i
        }
    }
    formatted := float64(size) / float64(d)
    return fmt.Sprintf("%.1f%s",formatted,sizes[index])
}
