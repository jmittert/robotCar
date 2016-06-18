package main

import (
  "fmt"
  "os"
  "io"
  "net"
  "io/ioutil"
  "github.com/BurntSushi/toml"
  "github.com/jmittert/xb360ctrl"
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
  xb360ctrl.DebugModeOn()
//  C.wiringPiSetup()
  //C.pinMode(C.A1, C.OUTPUT)
  //C.pinMode(C.A2, C.OUTPUT)
  //C.pinMode(C.B1, C.OUTPUT)
  //C.pinMode(C.B2, C.OUTPUT)

  var config Config
  f, err := ioutil.ReadFile("carrc");
  checkError(err)
  if _, err := toml.Decode(string(f), &config); err != nil {
    checkError(err)
  }
  service := config.ServerAddr

  bytes := make([]byte, 8)
  var xbState xb360ctrl.Xbc_state
  for {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    if err != nil {
      continue
    }
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
      continue
    }
    for {
      count, err := conn.Read(bytes)
      if count == 0 && err == io.EOF {
        // On EOF, disconnect and look for another connection
        fmt.Println("EOF!")
        conn.Close()
        break;
      }
      if count != 8 {
        fmt.Println("Got ", count, "/8 bytes")
        conn.Close()
        break;
      }
      var e xb360ctrl.Xbc_event
//      fmt.Println(bytes)
      e.UnMarshalBinary(bytes)
      xb360ctrl.UpdateState(&e, &xbState)
      if xbState.RTrigger > 5 {
        C.digitalWrite (C.A1, C.HIGH)
        C.digitalWrite (C.A2, C.LOW)
        C.digitalWrite (C.B1, C.HIGH)
        C.digitalWrite (C.B2, C.LOW)
      } else {
        C.digitalWrite (C.A1, C.LOW)
        C.digitalWrite (C.A2, C.LOW)
        C.digitalWrite (C.B1, C.LOW)
        C.digitalWrite (C.B2, C.LOW)
      }
//      fmt.Println(e)
 //     fmt.Println(xbState)
    }
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}

