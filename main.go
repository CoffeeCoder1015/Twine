package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
    keys *CustomKeyMap
    db *Debouncer
}

func initModel() model{
   return model{
        focus: 0,
        inputs: InitInput(),
        results: InitResults(),
        keys: newKeyMap(),
        db: NewDebouncer(time.Millisecond*150),
    } 
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return textinput.Blink
}

// Called to render UI
func (m model) View() string{
    header := TitleStyle.Render("🌳 Twine 🎄")
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
    case DebounceLoading:
        m.results.list.Title = "Loading..."
        return m, nil
    case DebounceRefresh:
        m.results.UpdateList(false)
        m.results, _ = m.results.Update(msg)
        n := len( m.results.twine.flatCache )
        if n > 300_000{
            m.db.debounceDelay = time.Millisecond*200
        }else if n > 150_000{
            m.db.debounceDelay = time.Millisecond*150
        }else if n > 50_000{
            m.db.debounceDelay = time.Millisecond*70
        }else{
            m.db.debounceDelay = time.Millisecond*10
        }
        return m,nil
    case tea.KeyMsg:
        switch{
        case key.Matches(msg,m.keys.SwitchPanel):
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
        }
        switch msg.Type {
        case tea.KeyCtrlC:
            return m, tea.Quit
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
            path := m.inputs.inputs[0].Value()
            m.results.twine.directory = path
            m.results.twine.filter = m.inputs.GetFilter()
            m.db.ScheduleUpdate(func() {
                m.results.UpdateList(false)
            })
        }
        if passThrough && !m.db.activeWorker {
            m.results, _ = m.results.Update(msg)
        }
    }else{
        // while focused to the results panel
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch {
            case key.Matches(msg,m.keys.JumpToSelected):
                // jump to item
                selected_index := m.results.index
                notInFilter := m.results.list.FilterState() == list.Unfiltered
                if notInFilter && 0 <= selected_index && int(selected_index) < len(m.results.twine.flatCache) {
                    selected := m.results.twine.flatCache[selected_index]
                    selected_dir := formatPath(selected.path)
                    if selected.IsDir(){
                        selected_dir = filepath.Join(selected.path,selected.Name())
                    }
                    m.refreshRootDirectory(selected_dir)
                }
            case key.Matches(msg,m.keys.GoToParentDir):
                // go up in directory
                current := m.inputs.inputs[0].Value()
                after := filepath.Dir(current)
                m.refreshRootDirectory(after)
            case key.Matches(msg,m.keys.LaunchDefault):
                // launch default app
                selected_index := m.results.index
                notInFilter := m.results.list.FilterState() == list.Unfiltered
                if notInFilter && 0 <= selected_index && int(selected_index) < len(m.results.twine.flatCache) {
                    selected := m.results.twine.flatCache[selected_index]
                    selected_dir := filepath.Join(selected.path,selected.Name())
                    launchDefaultApp(selected_dir)
                }
            case key.Matches(msg,m.keys.WriteResult):
                m.results.twine.writeResult(m.inputs.SearchPattern())
            case key.Matches(msg,m.keys.Refresh):
                m.results.UpdateList(true)
            }
        }
        m.results, cmd = m.results.Update(msg)
    }
    return m,cmd
}

func (m *model) refreshRootDirectory(selected_dir string){
    m.results.twine.directory = selected_dir
    m.inputs.inputs[0].SetValue(selected_dir)
    m.inputs.inputs[0].SetCursor(len(selected_dir))
    m.results.twine.directory = m.inputs.inputs[0].Value()
    m.results.UpdateList(false)
}

func launchDefaultApp(path string){
    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "windows":
        // On Windows, use "cmd /C start" to launch default app
        cmd = exec.Command("cmd", "/C", "start", path)
    case "darwin": // macOS
        // On macOS, use "open"
        cmd = exec.Command("open", path)
    case "linux":
        // On Linux, use "xdg-open"
        cmd = exec.Command("xdg-open", path)
    default:
        fmt.Println( fmt.Errorf("unsupported operating system: %s", runtime.GOOS) )
    }

    // Use cmd.Start() to launch the application without waiting for it to close
    err := cmd.Start()
    if err != nil {
        fmt.Println( fmt.Errorf("failed to launch default app for '%s' on %s: %w", path, runtime.GOOS, err) )
    }
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
    m := initModel()
    p := tea.NewProgram(m,tea.WithAltScreen())
    m.db.d = p
    if _, err := p.Run(); err != nil{
        fmt.Printf("error has occured: %v",err)
        os.Exit(1)
    }
}
