package main

import (
  "fmt"
  "os/exec"
)


func returnOutput(c exec.Cmd) ([]byte, error){
  return c.CombinedOutput
}

func main(){
  c := exec.Command("sysct", "-n", "machdep.cpu.brand_string")
  fmt.Println("\n")
  err := c.Run()
  out, err2 := returnOutput(c)
  if err != nil {
    fmt.Println("Error:", err)
  } else if err2 != nil {
    fmt.Println("Error2: ", err2)
  }
}
