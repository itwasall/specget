package main

import (
  "fmt"
  "os/exec"
)

func main(){
  c := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
  a, err := c.Output()
  if err != nil {
    fmt.Println("Error: ", err)
  } else {
    fmt.Println(string(a[:]))
  }
  
  fmt.Println(c.String())
}
