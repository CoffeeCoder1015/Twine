package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/bubbles/list"
)

var (
    sizes = []string{"b","Kb","Mb","Gb"}
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
    item
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

func (t Twine) Query() []list.Item{
    t.Search()
    cache := t.cache[t.filter.directory]
    len := min(10000,len(cache))
    cache = cache[:len]
    r := make( []list.Item, len )
    for i := range cache{
        r[i] = cache[i].item
    }
    return r
}

func (t Twine) SmartQuery(index int64) []list.Item{
    t.Search()
    cache := t.cache[t.filter.directory]

    upper := min( index+500,int64( len(cache)-1 ) )
    lower := max( index-500,0 )
    cache = cache[lower:upper]

    r := make( []list.Item, len(cache) )
    for i := range cache{
        r[i] = cache[i].item
    }
    return r
}

func (t Twine) Search(){
    _,c := t.cache[t.filter.directory]
    if c {
        return
    }
    
    results := make([]resultEntry,0)
    queue := []string{t.filter.directory}
    for 0 < len(queue){
        path := queue[0]
        queue = queue[1:]
        de, _ := os.ReadDir(path)
        for i := range de{
            e := de[i]
            r := resultEntry{path: path,DirEntry: e}
            r.formatInfo()
            results = append(results, r)
            if e.IsDir() {
                queue = append(queue, filepath.Join(path,e.Name())) 
            }
        } 
    }
    t.cache[t.filter.directory] = results
}

func (entry *resultEntry) formatInfo(){
    info, _ := entry.Info()
    icon := "📄"
    if entry.IsDir(){
        icon = "📁"
    }
    title := fmt.Sprintf("%s %s %s",entry.Name(),entry.path,icon)
    desc := fmt.Sprintf("%s %s %s",formatSize(info.Size()),info.ModTime().Format("2006-01-02 15:04:05"),info.Mode())
    entry.title = title
    entry.desc = desc
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
