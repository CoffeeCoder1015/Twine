package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ( 
    sizes = []string{"b","Kb","Mb","Gb"}
    docStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
)

type item struct{
    title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type ResultsList struct{
    list list.Model
    filter queryFilterPattern
}

func InitResults() ResultsList{
    delegate := list.NewDefaultDelegate()
    l := list.New(_tempList("."),delegate,0,0)
    l.Title = "Results"
    l.KeyMap.Quit.Unbind()
    l.KeyMap.ForceQuit.Unbind()

    return ResultsList{
        filter: queryFilterPattern{
            directory: ".",
        },
        list: l,
    }
}

func (m ResultsList) Init() tea.Cmd{
    return nil
}

func (m ResultsList) Update(msg tea.Msg) (ResultsList,tea.Cmd){
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        h, v := docStyle.GetFrameSize()
        m.list.SetSize(msg.Width-h, msg.Height-v-12)
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m ResultsList) View() string{
    return docStyle.Render(m.list.View())
}

func (m* ResultsList) UpdateList(){
    m.list.SetItems(_tempList(m.filter.directory))
}

func _tempList(dir string) []list.Item{
    entries, _ := os.ReadDir(dir)
    items :=  make([]list.Item,len(entries))
    for i, e := range entries{
        items[i] = formatInfo(e)
    } 
    return items
}

func formatInfo(entry os.DirEntry) item{
    info, _ := entry.Info()
    icon := "📄"
    if entry.IsDir(){
        icon = "📁"
    }
    title := fmt.Sprintf("%s %s",entry.Name(),icon)
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
