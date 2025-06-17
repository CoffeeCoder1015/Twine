package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ( 
    docStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
)

type CustomKeyMap struct{
    SwitchPanel key.Binding
    JumpToSelected key.Binding
    GoToParentDir key.Binding
    LaunchDefault key.Binding
    WriteResult key.Binding
    Refresh key.Binding
}

func newKeyMap() *CustomKeyMap{
    return &CustomKeyMap{
        SwitchPanel: key.NewBinding(
            key.WithKeys("shift+down","shift+up","alt+J","alt+K"),
            key.WithHelp("shift + ↑/↓", "Switch panel"),
            ),
        JumpToSelected: key.NewBinding(
            key.WithKeys("."),
            key.WithHelp(".","Jump to selected item"),
            ),
        GoToParentDir: key.NewBinding(
            key.WithKeys(","),
            key.WithHelp(",","Jump to parent directory"),
            ),
        LaunchDefault: key.NewBinding(
            key.WithKeys("enter"),
            key.WithHelp("enter","Launch default application"),
            ),
        WriteResult: key.NewBinding(
            key.WithKeys("w"),
            key.WithHelp("w","write result to disk"),
            ),
        Refresh: key.NewBinding(
            key.WithKeys("R"),
            key.WithHelp("shift + r","refresh catelogue"),
            ),
    }
}

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
    sliceLength int64
}

func InitResults() ResultsList{
    t := InitTwine()
    delegate := list.NewDefaultDelegate()
    t.Search(false)
    t.flattenTree()
    l := list.New(t.SmartQuery(0,1000),delegate,0,0)
    l.Title = "Results"
    l.KeyMap.Quit.Unbind()
    l.KeyMap.ForceQuit.Unbind()

    keymap := newKeyMap()
    l.AdditionalFullHelpKeys = func() []key.Binding {
        return []key.Binding{
            keymap.SwitchPanel,
            keymap.JumpToSelected,
            keymap.GoToParentDir,
            keymap.LaunchDefault,
            keymap.WriteResult,
            keymap.Refresh,
        }
    }

    l.AdditionalShortHelpKeys = func() []key.Binding {
        return []key.Binding{
            keymap.WriteResult,
            keymap.SwitchPanel,
        }
    }


    return ResultsList{
        twine: t,
        list: l,
        sliceLength: 1000,
    }
}

func (m ResultsList) Init() tea.Cmd{
    return nil
}

func (m ResultsList) Update(msg tea.Msg) (ResultsList,tea.Cmd){
    previous := m.index
    max_len := int64(len(m.twine.flatCache))

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        h, v := docStyle.GetFrameSize()
        m.list.SetSize(msg.Width-h, msg.Height-v-20)
        m.sliceLength = int64(m.list.Paginator.PerPage)*10 
    case tea.KeyMsg:
        // capturing keystrokes and updating global index
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

    // manually track global index because there are too many files
    if m.index < 0 {
        m.index = previous
    }else if m.index > max_len-1{
        m.index = max_len-1
        if m.list.Paginator.Page+2 == m.list.Paginator.TotalPages{
            m.index = max_len-1
        }
    }
    if m.index != previous {
        m.list.SetItems(m.twine.SmartQuery(m.index,m.sliceLength))
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)


    // Reset cursor position on new slice of items 
    diff := ( m.list.Index() )-int(m.index%m.sliceLength)
    action := m.list.CursorUp 
    if diff < 0{
        action = m.list.CursorDown
        diff = -diff
    }
    if m.list.FilterState() == list.Unfiltered {
        for range diff{
            action()
        }
    }

    m.list.Title =  fmt.Sprintf("Results %d/%d",m.index,max_len)

    return m, cmd
}

func (m ResultsList) View() string{
    return docStyle.Render(m.list.View())
}

func (m* ResultsList) UpdateList(refresh bool){
    m.index = 0
    m.twine.Search(refresh)
    m.twine.directory = formatPath(m.twine.directory)
    m.twine.flattenTree()
    r := m.twine.SmartQuery(m.index,m.sliceLength)
    m.list.SetItems(r)
}
