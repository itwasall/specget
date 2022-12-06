package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "golang.org/x/exp/slices"
)


func main(){
  files, err := ioutil.ReadDir("/")
  if err != nil {
    log.Fatal(err)
  }

  var filelist []string
  fmt.Println("FIRST PASS")
  for _, f := range files {
    filelist = append(filelist, f.Name())
  }
  fmt.Println("SECOND PASS")
  if slices.Contains(filelist, "bin"){
    fmt.Println("/bin/ found")
  }
  }

