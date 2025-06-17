package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

var (
    sizes = []string{"b","Kb","Mb","Gb"}
)

type filterFunc func(e os.DirEntry) bool

type resultEntry struct{
    fs.DirEntry
    path string 
    item
}

type cacheNode struct{
    r []resultEntry
    subdir []string
}

type Twine struct{
    directory string
    filter []filterFunc
    cache map[string]cacheNode
    flatCache []resultEntry
}

func InitTwine() Twine{
    t := Twine{
        directory:".",
        cache: make(map[string]cacheNode),
    }
    return t
}

func (t Twine) SmartQuery(index , width int64) []list.Item{
    t.Search(false)
    cache := t.flatCache

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

func (t Twine) Search(refresh bool){
    _,c := t.cache[t.directory]
    if !refresh{
        if c {
            return
        }
    }
    
    queue := []string{t.directory}
    cNode := make(chan cacheNode,3000)
    for 0 < len(queue){
        l := len( queue )
        for i := range l{
            go func ()  {
                results := make([]resultEntry,0)
                subdir := make([]string,0)
                path := queue[i]
                path = filepath.Clean(path) + "\\"
                de, _ := os.ReadDir(path)
                for i := range de{
                    e := de[i]
                    r := resultEntry{path: path,DirEntry: e}
                    r.formatInfo()
                    results = append(results, r)
                    if e.IsDir() {
                        next_path := filepath.Join(path,e.Name()) 
                        next_path = filepath.Clean(next_path) + "\\"
                        subdir = append(subdir, next_path)
                    }
                } 
                cNode <- cacheNode{r: results, subdir: subdir}
            }()
        }
        for i := range l{
            recv := <-cNode
            queue[i] = filepath.Clean(queue[i]) + "\\"
            t.cache[queue[i]] = recv
            queue = append(queue, recv.subdir...)
        }
        queue = queue[l:]
    }
    close(cNode)
}

func (t Twine) SearchSingle(){
    _,c := t.cache[t.directory]
    if c {
        return
    }
    
    queue := []string{t.directory}
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
        t.cache[path] = cacheNode{r: results, subdir: subdir}
    }
}

func (t *Twine) flattenTree() {
    r := make([]resultEntry,0)
    queue := []string{t.directory}
    chanList := make([]chan resultEntry,0)
    for 0 < len(queue){
        current := queue[0]
        queue = queue[1:]

        merge := t.cache[current]
        chanList = append(chanList, make(chan resultEntry,len(merge.r)))
        go func(index int){
            for _,item := range merge.r{
                success := true
                for _, ff := range t.filter{
                    success = ff(item)
                    if !success {
                    break 
                    }
                }
                if !success{
                    continue
                }
                chanList[index] <- item
            }
            close(chanList[index])
        }(len(chanList)-1)
        queue = append(queue, merge.subdir...)
    }
    for _, c := range chanList{
        for {
            v, ok := <- c
            if ok{
                r = append(r, v)
            }else{
                break
            }
        }
    }
    t.flatCache = r
}


func (t *Twine) flattenTreeSingle() {
    r := make([]resultEntry,0)
    queue := []string{t.directory}
    for 0 < len(queue){
        current := queue[0]
        queue = queue[1:]
        merge := t.cache[current]
        for _,item := range merge.r{
            success := true
            for _, ff := range t.filter{
                success = ff(item)
                if !success {
                   break 
                }
            }
            if !success{
                continue
            }
            r = append(r, item)
        }
        queue = append(queue, merge.subdir...)
    }
    t.flatCache = r 
}

func (t Twine) writeResult(header string){
    file, err := os.Create("result.twine.log")
    if err != nil {
       fmt.Println(err) 
        return
    }
    logString := header 
    for _, v:= range t.flatCache{
        icon := "📄 file"
        if v.IsDir(){
            icon = "📁 dir"
        }
        logString+=fmt.Sprintf("%s %s %s %s\n",icon,v.Name(),v.path,v.desc)
    }
    file.WriteString(logString)
    file.Close()
}


func (entry *resultEntry) formatInfo(){
    info, _ := entry.Info()
    icon := "📄"
    if entry.IsDir(){
        icon = "📁"
    }
    title := fmt.Sprintf("%s %s   \n%s   ",icon,entry.Name(),entry.path)
    desc := fmt.Sprintf("%s %s %s  ",formatSize(info.Size()),info.ModTime().Format("2006\\01\\02 15:04:05"),info.Mode())
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
