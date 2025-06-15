package main

import (
	"fmt"
	"os"

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
    inputs InputModel
    results ResultsList
}

func initModel() model{
   return model{
        inputs: InitInput(),
        results: InitResults(),
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
    s += m.inputs.View() + "\n"
    s += m.results.View()
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
        }
    }
    m.inputs, cmd = m.inputs.Update(msg)

    if m.inputs.validCount == len(m.inputs.inputs){
    }
    m.results.filter.directory = m.inputs.inputs[0].Value()
    newResults, resultsCmd := m.results.Update(msg)
    m.results = newResults
    cmd = tea.Batch(resultsCmd,cmd)

    return m,cmd
}


func main() {
    p := tea.NewProgram(initModel())
    if _, err := p.Run(); err != nil{
        fmt.Printf("error has occured: %v",err)
        os.Exit(1)
    }
}
