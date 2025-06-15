package main

import (
	"os"
	"regexp"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ( 
    docStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
)

type queryFilterPattern struct{
    directory string
    name *regexp.Regexp
    fileSize string
    mode *regexp.Regexp
    date string
    DirFile string
}

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
    l := list.New(_tempList("."),list.NewDefaultDelegate(),0,0)
    l.Title = "Results"
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
        items[i] = item{
            title: e.Name(),
            desc: strconv.FormatBool(e.IsDir()),
        }
    } 
    return items
}
