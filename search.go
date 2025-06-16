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

type cacheNode struct{
    depth int64
    r []resultEntry
    subdir []string
}

type Twine struct{
    filter queryFilterPattern
    cache map[string]cacheNode
}

func InitTwine() Twine{
    t := Twine{
        filter: queryFilterPattern{
            directory:".",
        },
        cache: make(map[string]cacheNode),
    }
    return t
}

func (t Twine) SmartQuery(index , width int64) []list.Item{
    t.Search()
    cache := t.flattenTree()

    m := index/width
    upper := min( m*width+width,int64( len(cache) ) )
    lower := m*width
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
    
    queue := []string{t.filter.directory}
    for 0 < len(queue){
        results := make([]resultEntry,0)
        subdir := make([]string,0)
        path := queue[0]
        queue = queue[1:]
        de, _ := os.ReadDir(path)
        for i := range de{
            e := de[i]
            r := resultEntry{path: path,DirEntry: e}
            r.formatInfo()
            results = append(results, r)
            if e.IsDir() {
                next_path := filepath.Join(path,e.Name()) 
                queue = append(queue, next_path)
                subdir = append(subdir, next_path)
            }
        } 
        t.cache[path] = cacheNode{r: results, subdir: subdir,depth: int64(len(de))}
    }
    queue = []string{t.filter.directory}
    for i := 0; i < len(queue); i++{
        current := queue[i]
        merge := t.cache[current]
        queue = append(queue, merge.subdir...)
    }
    for i := len(queue)-1; 0 <= i; i--{
        current := queue[i]
        current_cache := t.cache[current]
        for _,v := range current_cache.subdir{
            current_cache.depth += t.cache[v].depth
        }
        t.cache[current] = current_cache
    }
}

func (t Twine) flattenTree() []resultEntry{
    r := make([]resultEntry,0)
    queue := []string{t.filter.directory}
    for 0 < len(queue){
        current := queue[0]
        queue = queue[1:]
        merge := t.cache[current]
        r = append(r, merge.r...)
        queue = append(queue, merge.subdir...)
    }
    return r
    
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
