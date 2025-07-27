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

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type ResultsList struct {
	list        list.Model
	index       int64
	sliceLength int64
	cache       []resultEntry
}

func InitResults() ResultsList {
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
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
			keymap.Quit,
		}
	}

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keymap.WriteResult,
			keymap.SwitchPanel,
		}
	}

	return ResultsList{
		list:        l,
		sliceLength: 1000,
	}
}

func (m ResultsList) Init() tea.Cmd {
	return nil
}

func (m ResultsList) Update(msg tea.Msg) (ResultsList, tea.Cmd) {
	previous := m.index
	max_len := int64(len(m.cache))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-20)
		m.sliceLength = int64(m.list.Paginator.PerPage) * 10
	case tea.KeyMsg:
		// capturing keystrokes and updating global index
		switch msg.String() {
		case "down", "j":
			m.index++
		case "up", "k":
			m.index--
		case "right", "l", "pgdown", "f", "d":
			m.index += int64(m.list.Paginator.PerPage)
		case "left", "h", "pgup", "b", "u":
			m.index -= int64(m.list.Paginator.PerPage)
		case "home", "g":
			m.index = 0
		case "end", "G":
			m.index = max_len - 1
		}
	}

	// manually track global index because there are too many files
	if m.index < 0 {
		m.index = previous
	} else if m.index > max_len-1 {
		m.index = max_len - 1
		if m.list.Paginator.Page+2 == m.list.Paginator.TotalPages {
			m.index = max_len - 1
		}
	}
	if m.index != previous {
		m.list.SetItems(m.getResultWindow(m.index, m.sliceLength))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// Reset cursor position on new slice of items
	diff := (m.list.Index()) - int(m.index%m.sliceLength)
	action := m.list.CursorUp
	if diff < 0 {
		action = m.list.CursorDown
		diff = -diff
	}
	if m.list.FilterState() == list.Unfiltered {
		for range diff {
			action()
		}
	}

	m.list.Title = fmt.Sprintf("Results %d/%d", m.index, max_len)

	return m, cmd
}

func (m ResultsList) View() string {
	return docStyle.Render(m.list.View())
}

func (m ResultsList) getResultWindow(index, width int64) []list.Item {
	cache := m.cache

	mult := index / width
	upper := min(mult*width+width, int64(len(cache)))
	lower := mult * width
	cache = cache[lower:upper]

	r := make([]list.Item, len(cache))
	for i := range cache {
		r[i] = cache[i].item
	}
	return r
}
func (m *ResultsList) UpdateList(newItems []resultEntry) {
	m.cache = newItems
	m.index = 0
	r := m.getResultWindow(m.index, m.sliceLength)
	m.list.SetItems(r)
}
