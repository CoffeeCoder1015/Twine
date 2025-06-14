package main

import (
	"fmt"
    tea "github.com/charmbracelet/bubbletea"
)


type model struct{

}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

// Called to render UI
func (m model) View() string{
}


// Handles state updates like key inputs
func (m model) Update(msg tea.Msg) (tea.Model, tea.Msg){

}

func main() {
    fmt.Println("Hello world")
}
