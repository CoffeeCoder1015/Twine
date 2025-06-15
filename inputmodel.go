package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
    ModelStyle = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder())

    promptFocusedStyle = lipgloss.NewStyle()
    promptBluredStyle = lipgloss.NewStyle().
        Italic(true).
        Foreground(lipgloss.Color( "240" ))
)


type InputModel struct {
    focus int
    inputs  []textinput.Model
}

func InitInput() InputModel{
    m := InputModel{
        inputs: make([]textinput.Model, 6),
    }

    var input textinput.Model
    for i := range m.inputs{
        input = textinput.New()
        input.Width = 40
        input.PromptStyle = promptBluredStyle
        input.PlaceholderStyle = input.PlaceholderStyle. 
            Italic(true)
        switch i{
        case 0:
            input.Prompt = "Search directory: "
            input.Placeholder = "<current directory>"
        case 1:
            // set as focus
            input.Focus()
            m.focus = i
            input.PromptStyle = promptFocusedStyle

            input.Prompt = "Match pattern: "
            input.Placeholder = ".* "
        case 2:
            input.Prompt = "File size: "
            input.Placeholder = "> 0b"
        case 3:
            input.Prompt = "File mode: "
            input.Placeholder = "-----"
        case 4:
            input.Prompt = "Mod time: "
            input.Placeholder = "< 1978-03-12 13:23"
        case 5:
            input.Prompt = "Directory/File: "
            input.Placeholder = "dir / file / all"
        }

        m.inputs[i] = input
    }
    return m
}

func (m InputModel) Init() tea.Cmd{
    return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (InputModel,tea.Cmd ){
    var cmd tea.Cmd
    switch msg := msg.(type){
    case tea.KeyMsg:
        switch msg.Type{
        case tea.KeyTab, tea.KeyShiftTab:
            s := msg.String()

            if s == "tab" {
                m.focus++
            }else{
                m.focus--
            }

            if m.focus < 0{
                m.focus = len(m.inputs)-1
            }else if m.focus >= len(m.inputs){
                m.focus = 0
            }

            cmds := make([]tea.Cmd, len(m.inputs))
            for i := range m.inputs{
                if i == m.focus{
                    m.inputs[i].PromptStyle = promptFocusedStyle
                    continue
                }
                m.inputs[i].Blur()
                m.inputs[i].PromptStyle = promptBluredStyle
            }
            cmds[m.focus] = m.inputs[m.focus].Focus()
            return m,tea.Batch(cmds...)
        }
    }

    cmd = m.UpdateInputs(msg)
    return m,cmd
}

func (m* InputModel) UpdateInputs(msg tea.Msg) tea.Cmd{
    cmds := make([]tea.Cmd,len(m.inputs))

    for i := range m.inputs{
        m.inputs[i],cmds[i] = m.inputs[i].Update(msg)
    }

    return tea.Batch(cmds...)
}

func (m InputModel) View() string{
    views := make([]string,len(m.inputs))
    for i,v := range m.inputs{
       views[i] = v.View()
    }

    return ModelStyle.Render(strings.Join(views,"\n"))
}
