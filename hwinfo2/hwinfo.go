package main

import (
  "fmt"
  "os"
  "os/exec"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"

)

var (
  titleStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("#505059")).
    Foreground(lipgloss.Color("#ffffff"))
  bodyStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#a0f0f0"))
  commandStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#30b0e0"))

)

type model struct {
  altscreen bool
  command string
  err error 
}

type commandFinishedMsg struct { err error }

func (m model) Init() tea.Cmd {
  return nil
}

func execCommand(m model) tea.Cmd {
  command := "neofetch"

  c := exec.Command(command)

  return tea.ExecProcess(c, func(err error) tea.Msg {
    return commandFinishedMsg{err}
  })
}

func openEditor() tea.Cmd {
  editor := "nvim"

  c := exec.Command(editor)

  return tea.ExecProcess(c, func(err error) tea.Msg {
    return commandFinishedMsg{err}
  })
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type){
  case tea.KeyMsg:
    switch msg.String(){
    case "q":
      return m, tea.Quit
    case "d":
      return m, execCommand(m)
    case "e":
      return m, openEditor()
    case "a":
      m.altscreen = !m.altscreen
      cmd := tea.EnterAltScreen
      if !m.altscreen{
        cmd = tea.ExitAltScreen
      }
      return m, cmd
    }
  case commandFinishedMsg:
    if msg.err != nil {
      m.err = msg.err
      return m, tea.Quit
    }
  }
  return m, nil
}

func (m model) View() string {

  title := titleStyle.Render("Please enter a command")


  bodyNeofetch := bodyStyle.Render("\n'd' for ") + commandStyle.Render("neofetch")
  bodyNvim := bodyStyle.Render("\n'e' for ") + commandStyle.Render("nvim") 
  bodyQuit := bodyStyle.Render("\n'q' to ") + commandStyle.Render("quit")

  if m.err != nil {
    return "Error: " + m.err.Error() + "\n"
  }
  return title+bodyNeofetch+bodyNvim+bodyQuit+"\n"

}

func main(){
  m := model{}

  if err:= tea.NewProgram(m).Start(); err != nil {
    fmt.Println("Error! ", err)
    os.Exit(1)
  }
}
