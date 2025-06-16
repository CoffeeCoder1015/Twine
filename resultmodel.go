package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ( 
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
    twine Twine
    index int64
}

func InitResults() ResultsList{
    t := InitTwine()
    delegate := list.NewDefaultDelegate()
    l := list.New(t.Query(),delegate,0,0)
    l.Title = "Results"
    l.KeyMap.Quit.Unbind()
    l.KeyMap.ForceQuit.Unbind()

    return ResultsList{
        twine: t,
        list: l,
    }
}

func (m ResultsList) Init() tea.Cmd{
    return nil
}

func (m ResultsList) Update(msg tea.Msg) (ResultsList,tea.Cmd){
    // previous := m.index
    max_len := int64(len(m.twine.cache[m.twine.filter.directory]))

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        h, v := docStyle.GetFrameSize()
        m.list.SetSize(msg.Width-h, msg.Height-v-12)
    case tea.KeyMsg:
        switch msg.String(){
        case "down", "j":
            m.index++
        case "up", "k":
            m.index--
        case "right","l", "pgdown", "f", "d":
            m.index += int64(m.list.Paginator.PerPage) 
        case "left","h", "pgup", "b" ,"u":
            m.index -= int64(m.list.Paginator.PerPage) 
        case "home", "g":
            m.index = 0
        case "end", "G":
            m.index = max_len-1
        }
    }
    m.list.SetItems(m.twine.SmartQuery(m.index))        

    m.list.Title =  fmt.Sprintf("Results %d",m.index)
    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m ResultsList) View() string{
    return docStyle.Render(m.list.View())
}

func (m* ResultsList) UpdateList(){
    m.index = 0
    m.list.SetItems(m.twine.Query())
}
