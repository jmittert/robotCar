package main

import (
  "fmt"
  "os"
  "net"
  "github.com/BurntSushi/toml"
)
/*
#cgo LDFLAGS: -lwiringPi
#include <wiringPi.h>
*/
import "C"

type Config struct {
  controllerIP  string
}

func ReadConfig() Config {
  var config Config
  _, err := toml.DecodeFile("./carrc", &config);
  checkError(err)
  return config
}

func main() {
  config := ReadConfig()
  service := config.controllerIP

  for {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    if err != nil {
      continue
    }
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
      continue
    }
    fmt.Println("Connected")

    for {
      char := make([]byte, 1)
      _, err = conn.Read(char)
      print(string(char))
      if err != nil {
        fmt.Println("\nConnection lost")
        break
      }
    }
  }
  os.Exit(0)
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}

