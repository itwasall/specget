package main

import (
  "fmt"
  "bytes"
  //"encoding/json"
  "os"
  "os/exec"
  tea "github.com/charmbracelet/bubbletea"
)

type model struct {
  execType string
}

func Shellout(command string) (string, string, error) {
  var stdout bytes.Buffer
  var stderr bytes.Buffer

  cmd := exec.Command("sh", "-c", command)

  cmd.Stdout = &stdout
  cmd.Stderr = &stderr

  err := cmd.Run()

  return stdout.String(), stderr.String(), err
}

func ShelloutWrapper(m model, command string) (tea.Model, tea.Cmd) {
  stdout, stderr, err := Shellout(command)
  fmt.Println("---stdout---")
  fmt.Println(stdout)
  fmt.Println("---stderr---")
  fmt.Println(stderr)
  fmt.Println("---err---")
  fmt.Println(err)
  return m, nil
}

func stdoutPipe(m model) (tea.Model, tea.Cmd) {
  //cmd := exec.Command("echo", "-n", `{"Name": "Bob", "Age": 32}`)
  cmd := exec.Command("echo", "penis haha")
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    fmt.Println("StdoutPipe Error!", err)
    return m, nil
  }
  if err := cmd.Start(); err != nil {
    fmt.Println("StdoutPipe Start Error!", err)
    return m, nil
  }
  /*
  if err := json.NewDecoder(stdout).Decode(&person); err != nil {
    fmt.Println("Json Decode Error!", err)
    return m, nil
  }
  */
  
  if err := cmd.Wait(); err != nil {
    fmt.Println("CMD Wait Error!", err)
    return m, nil
  }

  fmt.Println(&stdout)

  return m, nil
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
  switch msg := msg.(type){
  case tea.KeyMsg:
    switch msg.String(){
    case "q":
      return m, tea.Quit
    case "1":
      return stdoutPipe(m)
    case "2":
      return ShelloutWrapper(m, `ioreg -rc IOPCIDevice | grep model`)
    }

  }
  return m, nil
}

func (m model) View() string {
  var renderString string
  renderString += "Exec Test\n"
  return renderString
}

func main(){
  m := model{}

  if err := tea.NewProgram(m).Start(); err != nil {
    fmt.Println("Boot error: ", err)
    os.Exit(1)
  }
}
