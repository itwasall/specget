package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"

)

var (
  titleStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("#505059")).
    Foreground(lipgloss.Color("#ffffff"))
  bodyStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#e030e0"))
  quitStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#303030")).
    PaddingLeft(3)
  commandStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#30f0e0"))

  borderStyle = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("#f0f0f0"))

)

type model struct {
  altscreen bool
  command string
  commandPresent bool
  err error 
}

type commandFinishedMsg struct { err error }

func (m model) Init() tea.Cmd {
  return nil
}

func getGPU(m model) (tea.Model, tea.Cmd){
  c:= exec.Command("bash", "-c", "ioreg -rc IOPCIDevice | grep \"model\" | sed -n '1 p' | awk '{print $5, $6, $7}'")
  //c:= exec.Command("bash", "-c", "ioreg -rc IOPCIDevice | grep \"model\"")
  gpu_info, err := c.Output()
  if err != nil {
    fmt.Println("Error:", err)
    return m, nil
  }

  m.command = string(gpu_info[:])
  return m, nil

}

func getRAM(m model) (tea.Model, tea.Cmd) {
  //c := exec.Command("sysctl", "hw.memsize")
  c := exec.Command("bash", "-c", "sysctl hw.memsize | awk '{print $2/1024/1024/1024 \"GB\"}'")
  // TODO
  // Implement command piping and implement awk filtering
  // awk := exec.Command("awk", "'{print $2/1024/1024/1024 \"GB\"}'")
  ram_info, err := c.Output()
  if err != nil {
    fmt.Println("Error:", err)
    return m, nil
  }
  m.command = (string(ram_info[:]))
  return m, nil

}

func getCPU(m model) (tea.Model, tea.Cmd) {
  c := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
  cpu_info, err := c.Output()
  if err != nil {
    fmt.Println("Error: ", err)
    return m, nil
  }
  m.command = string(cpu_info[:])
  //return tea.ExecProcess(c, func(err error) tea.Msg {
  //  return commandFinishedMsg{err}
  //})
  return m, nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type){
  case tea.KeyMsg:
    switch msg.String(){
    case "q":
      return m, tea.Quit
    case "d":
      tea.Println("")
      m.commandPresent = true
      return getCPU(m)
    case "r":
      m.commandPresent = true
      return getRAM(m)
    case "g":
      m.commandPresent = true
      return getGPU(m)
    case "a":
      m.altscreen = !m.altscreen
      cmd := tea.EnterAltScreen
      if !m.altscreen{
        cmd = tea.ExitAltScreen
      }
      return m, cmd
    case "x":
      return m, nil
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
  var renderString string
  if !m.commandPresent {
    renderString = titleStyle.Render("Please enter a command")

    renderString += bodyStyle.Render("\n'd' for ") + commandStyle.Render("cpu")
    renderString += bodyStyle.Render("\n'r' for ") + commandStyle.Render("ram") 
    renderString += bodyStyle.Render("\n'g' for ") + commandStyle.Render("gpu") 
  } else {
    renderString = "\n"
    switch {
    case strings.Contains(m.command, "Radeon"):
      renderString += bodyStyle.Render("\nGPU is: ") + commandStyle.Render(m.command)
    case strings.Contains(m.command, "Nvidia"):
      renderString += bodyStyle.Render("\nGPU is: ") + commandStyle.Render(m.command)
    case strings.Contains(m.command, "GB"):
      renderString += bodyStyle.Render("\nRAM is: ") + commandStyle.Render(m.command) 
    case strings.Contains(m.command, "Intel(R)"):
      renderString += bodyStyle.Render("\nCPU is: ") + commandStyle.Render(m.command)
    }
    m.commandPresent = !m.commandPresent
  }
  /*
  if m.command.Contains("GB") {
    m.commandPresent = !m.commandPresent
  } else if m.command.Contains("Intel") {
    renderString = "\n"
    m.commandPresent = !m.command
  }
  */
  renderString += quitStyle.Render("\n\n'q' to quit")

  if m.err != nil {
    return "Error: " + m.err.Error() + "\n"
  }
  return borderStyle.Render(renderString)

}

func main(){
  m := model{commandPresent: false}

  if err:= tea.NewProgram(m).Start(); err != nil {
    fmt.Println("Error! ", err)
    os.Exit(1)
  }
}
