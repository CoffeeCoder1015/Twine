package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
    HeaderStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#f0ece1")).
        Background(lipgloss.Color("#68c238"))
    TitleStyle = HeaderStyle. 
        Width(50).
        Align(lipgloss.Center).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#33d483"))
)

type model struct{
    regexInput textinput.Model
    rootInput textinput.Model
}

func initModel() model{
    regex := textinput.New()
    regex.Focus()
    regex.Placeholder = ".* "
    regex.PlaceholderStyle.Italic(true)
    regex.Prompt = "Match pattern:"
    regex.Width = 40

    root := textinput.New()
    root.PlaceholderStyle.Italic(true)
    root.Prompt = "Search directory:"
    root.Placeholder = "<current directory>"
    root.Width = 40

    return model{
        regexInput: regex,
        rootInput: root,
    }
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return textinput.Blink
}

// Called to render UI
func (m model) View() string{
    header := TitleStyle.Render(">> Twine <<")
    s := header+"\n"
    s+=m.regexInput.View()+"\n"
    s+=m.rootInput.View()+"\n"

    input_length := len(m.regexInput.Value())
    s+=strconv.Itoa(input_length)
    return s
}


// Handles state updates like key inputs
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
    var cmd tea.Cmd
   
    switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
        case tea.KeyTab, tea.KeyShiftTab:
            if m.regexInput.Focused() {
                m.regexInput.Blur() 
                m.rootInput.Focus()
            }else{
                m.regexInput.Focus() 
                m.rootInput.Blur()
            }
		}
	}

    updatedRegex, regxCmd := m.regexInput.Update(msg)
    updatedRoot, rootCmd := m.rootInput.Update(msg)
    m.regexInput = updatedRegex
    m.rootInput = updatedRoot
    cmd = tea.Batch(regxCmd,rootCmd)
    return m,cmd
}


func main() {
    p := tea.NewProgram(initModel())
    if _, err := p.Run(); err != nil{
        fmt.Printf("error has occured: %v",err)
        os.Exit(1)
    }
}
