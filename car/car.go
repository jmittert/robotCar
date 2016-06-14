package main

import (
  "fmt"
  "os"
  "net"
  "bufio"
  "io/ioutil"
  "github.com/BurntSushi/toml"
)
/*
#cgo LDFLAGS: -lwiringPi
#include <wiringPi.h>
*/
import "C"

type Config struct {
  ServerAddr string
}

func main() {
  var config Config
  f, err := ioutil.ReadFile("carrc");
  checkError(err)
  if _, err := toml.Decode(string(f), &config); err != nil {
    checkError(err)
  }
  service := config.ServerAddr

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

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
      fmt.Println(scanner.Text())
    }
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}

