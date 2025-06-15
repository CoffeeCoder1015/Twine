package main

import (
	"regexp"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ( 
    docStyle = lipgloss.NewStyle().Margin(1, 2) 
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
    return ResultsList{
        filter: queryFilterPattern{},
    }
}

func (m ResultsList) Init() tea.Cmd{
    return nil
}

func (m ResultsList) Update(msg tea.Msg) (tea.Model,tea.Cmd){
    	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ResultsList) View() string{
    return docStyle.Render(m.list.View())
}
