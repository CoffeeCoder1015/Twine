package main

import (
	"fmt"
	"os"
	"regexp"
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

    valueCorrectStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#68c238")).
        Italic(true)

    valueIncorrectStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#ff3b14")).
        Italic(true)

    fileSizeMatcher = regexp.MustCompile("^\\d+(?:[kmg]i?b|b)$")
    dateMatcher = regexp.MustCompile("^(?:\\d{4}-\\d{2}-\\d{2}|today)(?: \\d{2}:\\d{2}:\\d{2})?$")
)


type InputModel struct {
    focus int
    inputs  []textinput.Model
    nameRgex *regexp.Regexp
    modeRgex *regexp.Regexp
    validCount int
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
            input.Placeholder = "<directory>"
            input.TextStyle = valueCorrectStyle
            wd, err :=  os.Getwd()
            if err == nil {
                input.SetValue(wd)
            }else{
                fmt.Println(err)
            }
        case 1:
            // set as focus
            input.Focus()
            m.focus = i
            input.PromptStyle = promptFocusedStyle

            input.Prompt = "Match pattern: "
            input.Placeholder = ".* "
        case 2:
            input.Prompt = "File size: "
            input.Placeholder = "0b-"
        case 3:
            input.Prompt = "File mode: "
            input.Placeholder = ".*"
        case 4:
            input.Prompt = "Mod time: "
            input.Placeholder = "1970 00:00:00-today"
        case 5:
            input.Prompt = "Directory/File: "
            input.Placeholder = "dir|file|all"
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

            if m.inputs[m.focus].Value() == ""{
                m.autoFill()
            }

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
    m.ProcessInputs()

    return tea.Batch(cmds...)
}

func (m* InputModel) ProcessInputs() {
    m.validCount = 0
    for i := range m.inputs{
        valid := false
        input := m.inputs[i]
        switch i {
        // Valid search directory
        case 0: 
            _, err := os.Stat(input.Value())
            valid = err == nil
        // Valid regex
        case 1:
            regexObject, err := regexp.Compile(input.Value())
            valid = err == nil
            if err == nil{
                m.nameRgex = regexObject
            }else{
                m.nameRgex = nil
            }
        // Valid file size range
        case 2:
            base := input.Value()
            if strings.Contains(base,"-"){
                split := strings.Split(base,"-")
                if split[0] == "" && split[1] == ""{
                    // faulty range
                    valid = false
                }else{
                    lowerBound := fileSizeMatcher.MatchString(split[0]) || split[0] == ""
                    upperBound := fileSizeMatcher.MatchString(split[1]) || split[1] == ""
                    valid = upperBound && lowerBound
                }
            }else{
                // no range provided
                valid = false
            }
        case 3:
            regexObject, err := regexp.Compile(input.Value())
            valid = err == nil
            if err == nil{
                m.modeRgex = regexObject
            }else{
                m.modeRgex = nil
            }
        case 4:
            base := input.Value()
            if strings.Contains(base,"-"){
                split := strings.Split(base,"-")
                if split[0] == "" && split[1] == ""{
                    // faulty range
                    valid = false
                }else{
                    lowerBound := dateMatcher.MatchString(split[0]) || split[0] == ""
                    upperBound := dateMatcher.MatchString(split[1]) || split[1] == ""
                    valid = upperBound && lowerBound
                }
            }else{
                // no range provided
                valid = false
            }
        case 5:
            switch input.Value(){
            case "all", "file", "dir":
                valid = true
            default:
                valid = false
            }
        } 

        if valid{
            m.inputs[i].TextStyle = valueCorrectStyle
            m.validCount++
        }else{
            m.inputs[i].TextStyle = valueIncorrectStyle
        }
    }

}

func (m* InputModel) autoFill(){
    input := &m.inputs[m.focus]
    switch m.focus{
    case 0:
        wd, err :=  os.Getwd()
        if err == nil {
            input.SetValue(wd)
        }else{
            fmt.Println(err)
        }
    case 1:
        input.SetValue(".*")
    case 2:
        input.SetValue("0b-")
    case 3:
        input.SetValue(".*")
    case 4:
        input.SetValue("-today")
    case 5:
        input.SetValue("all")
    }
}

func (m InputModel) View() string{
    views := make([]string,len(m.inputs))
    for i,v := range m.inputs{
       views[i] = v.View()
    }

    return ModelStyle.Render(strings.Join(views,"\n"))
}
