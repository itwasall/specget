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
    Foreground(lipgloss.Color("#404040"))
  quitStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#a0a0a0")).
    PaddingLeft(3).
    PaddingRight(3)
  commandStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#0050c0"))

  borderStyle = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("#000000")).
    PaddingLeft(3).
    PaddingRight(3)


)

type model struct {
  altscreen bool
  command string
  commandPresent bool
  commandType string
  err error 
}

type commandFinishedMsg struct { err error }

func (m model) Init() tea.Cmd {
  return nil
}

func getHDD(m model) (tea.Model, tea.Cmd){
  c1 := exec.Command("bash", "-c", "diskutil info /dev/disk0 | grep \"Disk Size\" | awk '{print $3, $4, $5, $6}'")
  c2 := exec.Command("bash", "-c", "diskutil info /dev/disk1 | grep \"Disk Size\" | awk '{print $3, $4, $5, $6}'")
  c3 := exec.Command("bash", "-c", "diskutil info /dev/disk2 | grep \"Disk Size\" | awk '{print $3, $4, $5, $6}'")

  hdd_info, err := c1.Output()
  if err != nil {
    fmt.Println("Error!", err)
  }
  hdd_info2, err := c2.Output()
  if err != nil {
    fmt.Println("Error!", err)
  }
  hdd_info3, err := c3.Output()
  if err != nil {
    fmt.Println("Error!", err)
  }

  m.command = bodyStyle.Render("\n/dev/disk0: ") + commandStyle.Render(string(hdd_info[:])) + bodyStyle.Render("\n/dev/disk1: ")+ commandStyle.Render(string(hdd_info2[:])) + bodyStyle.Render("\n/dev/disk2: ") + commandStyle.Render(string(hdd_info3[:]))
  m.commandType = "HDD"
  return m, nil
}

func getGPU(m model) (tea.Model, tea.Cmd){
  c:= exec.Command("bash", "-c", "ioreg -rc IOPCIDevice | grep \"model\" | sed -n '1 p' | awk '{print $5, $6, $7, $8}'")
  //c:= exec.Command("bash", "-c", "ioreg -rc IOPCIDevice | grep \"model\"")
  gpu_info, err := c.Output()
  if err != nil {
    fmt.Println("Error:", err)
    return m, nil
  }

  m.command = string(gpu_info[:])
  m.commandType = "GPU"
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
  m.commandType = "RAM"
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
  m.commandType = "CPU"
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
    case "h":
      m.commandPresent = true
      return getHDD(m)
    case "c":
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

    renderString += "\n" + bodyStyle.Render("'c' for ") + commandStyle.Render("cpu")
    renderString += "\n" + bodyStyle.Render("'r' for ") + commandStyle.Render("ram") 
    renderString += "\n" + bodyStyle.Render("'g' for ") + commandStyle.Render("gpu") 
    renderString += "\n" + bodyStyle.Render("'h' for ") + commandStyle.Render("hdd")
    renderString += "\n" + quitStyle.Render("'q' to quit")
  } else {
    renderString = "\n"
    switch m.commandType {
    case "GPU":
      renderString += bodyStyle.Render("GPU is: ") + commandStyle.Render(m.command)
    case "RAM":
      renderString += bodyStyle.Render("RAM is: ") + commandStyle.Render(m.command)
    case "CPU":
      renderString += bodyStyle.Render("CPU is: ") + commandStyle.Render(m.command)
    case "HDD":
      renderString += bodyStyle.Render("Disk Size is: ") + commandStyle.Render(m.command)
    default:
      renderString += bodyStyle.Render("Fucking uhhhhhh") + commandStyle.Render(" Idk B")
    }
    m.commandPresent = !m.commandPresent
    renderString += "\n" + quitStyle.Render("[c]pu | [g]pu | [r]am | [h]dd | [q]uit")
  }
  /*
  if m.command.Contains("GB") {
    m.commandPresent = !m.commandPresent
  } else if m.command.Contains("Intel") {
    renderString = "\n"
    m.commandPresent = !m.command
  }
  */

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
