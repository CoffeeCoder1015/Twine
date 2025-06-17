package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
    ModelStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

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

    fileSizeMatcher = regexp.MustCompile("^\\d+(?:\\.\\d+)?(?:[kmg]i?b)$|^\\d+b$")
    dateMatcher = regexp.MustCompile("^(?:\\d{4}\\\\\\d{2}\\\\\\d{2}|today)(?: \\d{2}:\\d{2}:\\d{2})?$")
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
            input.Placeholder = "1970\\01\\12 00:00:00-today"
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
        // switching inputs
        case tea.KeyTab, tea.KeyShiftTab:
            s := msg.String()

            // autofill input if it is empty and user leaves focus of input
            if m.inputs[m.focus].Value() == ""{
                m.autoFill()
                m.ProcessInputs()
            }

            if s == "tab" {
                m.focus++
            }else{
                m.focus--
            }

            // make sure focus stays in bound
            if m.focus < 0{
                m.focus = len(m.inputs)-1
            }else if m.focus >= len(m.inputs){
                m.focus = 0
            }

            // setting focused input
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
            if strings.Count(base,"-") == 1{
                split := strings.Split(base,"-")
                lowerBound := fileSizeMatcher.MatchString(split[0]) || split[0] == ""
                upperBound := fileSizeMatcher.MatchString(split[1]) || split[1] == ""
                valid = upperBound && lowerBound
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
            if strings.Count(base,"-") == 1{
                split := strings.Split(base,"-")
                lowerBound := dateMatcher.MatchString(split[0]) || split[0] == ""
                upperBound := dateMatcher.MatchString(split[1]) || split[1] == ""
                valid = upperBound && lowerBound
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
        input.SetValue("-")
    case 3:
        input.SetValue(".*")
    case 4:
        input.SetValue("-")
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

func (m InputModel) SearchPattern() string{
    views := make([]string,len(m.inputs))
    for i,v := range m.inputs{
       views[i] = v.Prompt + " " + v.Value()
    }
    return ModelStyle.Render( strings.Join(views,"\n") )
}

func (m InputModel) GetFilter() []filterFunc{
    filter := []filterFunc{}
    if m.inputs[1].Value() != ".*" {
        filter = append(filter, func(e resultEntry) bool {
            return m.nameRgex.MatchString(e.Name())
        })
    }

    size_str := m.inputs[2].Value()
    if size_str != "-"{
        SizeBound := strings.Split(size_str, "-")
        lower := len( SizeBound[0] ) > 0
        upper := len( SizeBound[1] ) > 1
        if lower && !upper {
            filter = append(filter, func(e resultEntry) bool {
                info, _ := e.Info() 
                return strToBytes(SizeBound[0]) < info.Size() 
            })    
        }else if !lower && upper{
            filter = append(filter, func(e resultEntry) bool {
                info, _ := e.Info() 
                return  info.Size() < strToBytes(SizeBound[1])
            })    
        }else if lower && upper{
            filter = append(filter, func(e resultEntry) bool {
                info, _ := e.Info() 
                return  strToBytes(SizeBound[0]) < info.Size() && info.Size() < strToBytes(SizeBound[1])
            })    
        }
    }

    if m.inputs[3].Value() != ".*" {
        filter = append(filter, func(e resultEntry) bool {
            info, _ := e.Info()
            return m.modeRgex.MatchString(info.Mode().String())
        })
    }

    time_str := m.inputs[4].Value()
    if time_str != "-"{
        TimeBound := strings.Split(time_str, "-")
        lower := len( TimeBound[0] ) > 0
        upper := len( TimeBound[1] ) > 1
        if lower && !upper {
            filter = append(filter, func(e resultEntry) bool {
                info, _ := e.Info() 
                return info.ModTime().After(strToTime(TimeBound[0],false))
            })    
        }else if !lower && upper{
            filter = append(filter, func(e resultEntry) bool {
                info, _ := e.Info() 
                return info.ModTime().Before(strToTime(TimeBound[1],true))
            })    
        }else if lower && upper{
            filter = append(filter, func(e resultEntry) bool {
                info, _ := e.Info() 
                lowerTime :=  info.ModTime().After(strToTime(TimeBound[0],false))
                upperTime :=  info.ModTime().Before(strToTime(TimeBound[1],true))
                return lowerTime && upperTime
            })    
        }
    }

    switch m.inputs[5].Value(){
    case "dir":
        filter = append(filter, func(e resultEntry) bool {
            return e.IsDir()
        })
    case "file":
        filter = append(filter, func(e resultEntry) bool {
            return !e.IsDir()
        })
    }
    return filter
}

// converts a size string e.g (120mb, 23.1gib, etc)
// into its corresponding integer size as bytes
func strToBytes(rawSize string) int64{
    unit, _ := regexp.Compile("[kmg]i?b|b$")
    unit_loc := unit.FindStringIndex(rawSize)[0]
    value := rawSize[:unit_loc]
    unit_str := rawSize[unit_loc:]
    
    size, _ := strconv.ParseFloat(value,64)
    switch unit_str{
    case "gb":
        size*=1000000000
    case "mb":
        size*=1000000
    case "kb":
        size*=1000
    case "gib":
        size*=(1<<30)
    case "mib":
        size*=(1<<20)
    case "kib":
        size*=(1<<10)
    }
    return int64(size)
}

func strToTime(rawTime string,isEndBound bool) time.Time{
    date := "2006\\01\\02"
    time_str := " 15:04:05"
    parse := date
    has_time := strings.Count(rawTime," ") == 1
    if has_time{
        parse += time_str
    }

    //replacing shorthand of 'today'
    now := time.Now()
    rawTime = strings.Replace(rawTime,"today",now.Format(date),1)

    t, _ := time.Parse(parse,rawTime)
    if isEndBound && !has_time {
        t = t.AddDate(0,0,1)
    }
    return t
}
