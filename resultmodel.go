package main

import (
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
    m.list.SetItems(m.twine.Query())
}

