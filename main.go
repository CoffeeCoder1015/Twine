package main

import (
	"fmt"
	"os"
	"path/filepath"

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
    focus int
    inputs InputModel
    results ResultsList
}

func initModel() model{
   return model{
        focus: 0,
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

    passThrough := true
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            return m, tea.Quit
        case tea.KeyShiftDown,tea.KeyShiftUp:
            if m.focus == 0{
                m.focus = 1
                ModelStyle = ModelStyle.Border(lipgloss.HiddenBorder())
                docStyle = docStyle.Border(lipgloss.NormalBorder())
            }else{
                m.focus = 0
                docStyle = docStyle.Border(lipgloss.HiddenBorder())
                ModelStyle = ModelStyle.Border(lipgloss.NormalBorder())
            }
            return m, cmd
        default:
            passThrough = false
        }
    }


    if m.focus == 0{
        // while focused to the search panel
        oldInput := make([]string,len(m.inputs.inputs))
        for i := range len(m.inputs.inputs){
            oldInput[i] = m.inputs.inputs[i].Value()
        }
        m.inputs, cmd = m.inputs.Update(msg)
        newInput := make([]string,len(m.inputs.inputs))
        for i := range len(m.inputs.inputs){
            newInput[i] = m.inputs.inputs[i].Value()
        }
        if len(m.inputs.inputs) == m.inputs.validCount && compareInput(oldInput,newInput){
            m.results.twine.directory = m.inputs.inputs[0].Value()
            m.results.twine.filter = m.inputs.GetFilter()
            m.results.UpdateList()
        }
        if passThrough {
            m.results, _ = m.results.Update(msg)
        }
    }else{
        // while focused to the results panel
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case ".":
                selected_index := m.results.index
                selected := m.results.twine.flatCache[selected_index]
                if selected.IsDir(){
                    selected_dir := filepath.Join(selected.path,selected.Name())
                    m.results.twine.directory = selected_dir
                    m.inputs.inputs[0].SetValue(selected_dir)
                    m.inputs.inputs[0].SetCursor(len(selected_dir))
                    m.results.twine.directory = m.inputs.inputs[0].Value()
                    m.results.UpdateList()
                }
            }
        }
        m.results, cmd = m.results.Update(msg)
    }
    return m,cmd
}

func compareInput(old_input, new_input []string) bool{
    for i,v := range old_input{
        if new_input[i] != v{
            return true
        }
    }
    return false
}


func main() {
    p := tea.NewProgram(initModel())
    if _, err := p.Run(); err != nil{
        fmt.Printf("error has occured: %v",err)
        os.Exit(1)
    }
}
