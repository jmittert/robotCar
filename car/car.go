package main

import (
  "fmt"
  "os"
  "net"
  "strings"
  "bufio"
  "io/ioutil"
  "github.com/BurntSushi/toml"
)
/*
#cgo LDFLAGS: -lwiringPi
#include <wiringPi.h>
int A1 = 0;
int A2 = 1;
int B1 = 3;
int B2 = 4;
*/
import "C"

type Config struct {
  ServerAddr string
}

func main() {
  C.wiringPiSetup()
  C.pinMode(C.A1, C.OUTPUT)
  C.pinMode(C.A2, C.OUTPUT)
  C.pinMode(C.B1, C.OUTPUT)
  C.pinMode(C.B2, C.OUTPUT)

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
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
      fmt.Println(scanner.Text())
      strs := strings.Split(scanner.Text(), " ")
      if strs[0] == "5" {
        if strs[1][0] != '-' {
          C.digitalWrite(C.A1, C.HIGH)
          C.digitalWrite(C.A2, C.LOW)
          C.digitalWrite(C.B1, C.HIGH)
          C.digitalWrite(C.B2, C.LOW)
        } else {
          C.digitalWrite(C.A1, C.LOW)
          C.digitalWrite(C.A2, C.LOW)
          C.digitalWrite(C.B1, C.LOW)
          C.digitalWrite(C.B2, C.LOW)
        }
      }
    }
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}

