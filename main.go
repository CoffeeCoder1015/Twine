package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


type model struct{
    regexInput textinput.Model
}

func initModel() model{
    regex := textinput.New()
    regex.Focus()
    regex.Placeholder = ".* "
    regex.PlaceholderStyle.Italic(true)
    regex.Prompt = "Search pattern:"
    regex.Width = 40

    return model{
        regexInput: regex,
    }
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return textinput.Blink
}

// Called to render UI
func (m model) View() string{
    s := ">> Twine <<\n"
    s+=m.regexInput.View()+"\n"

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
		}
	}

    m.regexInput, cmd = m.regexInput.Update(msg)
    return m,cmd
}

func main() {
    p := tea.NewProgram(initModel())
    if _, err := p.Run(); err != nil{
        fmt.Printf("error has occured: %v",err)
        os.Exit(1)
    }
}
