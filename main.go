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

func getRAM(m model) tea.Cmd {
  c := exec.Command("sysctl", "hw.memsize")
  // TODO
  // Implement command piping and implement awk filtering
  // awk := exec.Command("awk", "'{print $2/1024/1024/1024 \"GB\"}'")
  ram_info, err := c.Output()
  if err != nil {
    fmt.Println("Error:", err)
    return nil
  }
  fmt.Println(string(ram_info[:]))
  return nil

}

func getCPU(m model) tea.Cmd {
  c := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
  cpu_info, err := c.Output()
  if err != nil {
    fmt.Println("Error: ", err)
    return nil
  }
  fmt.Println(string(cpu_info[:]))
  //return tea.ExecProcess(c, func(err error) tea.Msg {
  //  return commandFinishedMsg{err}
  //})
  return nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type){
  case tea.KeyMsg:
    switch msg.String(){
    case "q":
      return m, tea.Quit
    case "d":
      tea.Println("")
      return m, getCPU(m)
    case "r":
      return m, getRAM(m)
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


  bodyCPU := bodyStyle.Render("\n'd' for ") + commandStyle.Render("cpu")
  bodyRAM := bodyStyle.Render("\n'r' for ") + commandStyle.Render("ram") 
  bodyQuit := bodyStyle.Render("\n'q' to ") + commandStyle.Render("quit")

  if m.err != nil {
    return "Error: " + m.err.Error() + "\n"
  }
  return title+bodyCPU+bodyRAM+bodyQuit+"\n"

}

func main(){
  m := model{}

  if err:= tea.NewProgram(m).Start(); err != nil {
    fmt.Println("Error! ", err)
    os.Exit(1)
  }
}
