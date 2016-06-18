package main

import (
  "fmt"
  "os"
  "io"
  "net"
  "io/ioutil"
  "github.com/BurntSushi/toml"
  xbc "github.com/jmittert/xb360ctrl"
  "flag"
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
  checkFlags()
  config := readConfig()
  serverAddr := config.ServerAddr
  bytes := make([]byte, 8)
  var xbState xbc.Xbc_state
  xbc.PrepState(&xbState)

  for {
    conn := connect(serverAddr)
    for {
      count, err := conn.Read(bytes)
      if count == 0 && err == io.EOF {
        // On EOF, disconnect and look for another connection
        fmt.Println("EOF!")
        conn.Close()
        break;
      } else if count != 8 {
        fmt.Println("Got ", count, "/8 bytes")
        conn.Close()
        break;
      }
      var e xbc.Xbc_event
      e.UnMarshalBinary(bytes)
      xbc.UpdateState(&e, &xbState)
      stateToHw(&xbState)
    }
  }
}

// Uses the current state of the controller to set the appropriate hw pins
func stateToHw(state *xbc.Xbc_state) {
  if state.RTrigger > -22767 || state.A{
    xbc.DEBUG("Fowards")
    C.digitalWrite (C.A1, C.HIGH)
    C.digitalWrite (C.A2, C.LOW)
    C.digitalWrite (C.B1, C.HIGH)
    C.digitalWrite (C.B2, C.LOW)
  } else if state.LTrigger > -22767 || state.B {
    xbc.DEBUG("Backwards")
    C.digitalWrite (C.A1, C.LOW)
    C.digitalWrite (C.A2, C.HIGH)
    C.digitalWrite (C.B1, C.LOW)
    C.digitalWrite (C.B2, C.HIGH)
  } else if state.DPadX == 32767 {
    xbc.DEBUG("RIGHT")
    C.digitalWrite (C.A1, C.LOW)
    C.digitalWrite (C.A2, C.LOW)
    C.digitalWrite (C.B1, C.LOW)
    C.digitalWrite (C.B2, C.HIGH)
  } else if state.DPadX == -32767 {
    xbc.DEBUG("LEFT")
    C.digitalWrite (C.A1, C.LOW)
    C.digitalWrite (C.A2, C.HIGH)
    C.digitalWrite (C.B1, C.LOW)
    C.digitalWrite (C.B2, C.LOW)
  } else {
    xbc.DEBUG("STOP")
    C.digitalWrite (C.A1, C.LOW)
    C.digitalWrite (C.A2, C.LOW)
    C.digitalWrite (C.B1, C.LOW)
    C.digitalWrite (C.B2, C.LOW)
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}

// checkFlags is a help function to read and handle the command line flags
func checkFlags() {
  dbgPtr := flag.Bool("dbg", false, "Set debug mode")
  hwPtr := flag.Bool("hw", true, "Enable talking with hardware")

  flag.Parse()
  if *dbgPtr {
    xbc.DebugModeOn()
  }

  if *hwPtr {
    C.wiringPiSetup()
    C.pinMode(C.A1, C.OUTPUT)
    C.pinMode(C.A2, C.OUTPUT)
    C.pinMode(C.B1, C.OUTPUT)
    C.pinMode(C.B2, C.OUTPUT)
  }

}

func readConfig() Config {
  var config Config
  f, err := ioutil.ReadFile("carrc");
  checkError(err)
  _, err = toml.Decode(string(f), &config)
  checkError(err)
  return config
}

func connect(addr string) net.Conn {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
    checkError(err)
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    checkError(err)
    return conn
}
