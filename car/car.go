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
#cgo LDFLAGS: -lwiringPi -lpthread
#include <wiringPi.h>
#include <softPwm.h>
int A1 = 0;
int A2 = 1;
int LPWM = 5;
int B1 = 3;
int B2 = 4;
int RPWM = 6;
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
  // Calculate pwm
  var basePwm int = 100
  var leftMod float32 = 1
  var rightMod float32 = 1
  var newPwm int
  if state.LStickY > 1000 {
    // >1000 -> Go right -> slow down right wheel
    leftMod -= float32(int(state.LStickY)/32768)
    if leftMod < 0 {
      leftMod = 0
    }
  } else if state.LStickY < 1000 {
    // <1000 -> Go left -> slow down left wheel
    rightMod -= float32(int(state.LStickY)/-32768)
    if rightMod < 0 {
      rightMod = 0
    }
  }

  if state.RTrigger > -22767 || state.A {
    xbc.DEBUG("Fowards")
    C.digitalWrite (C.A1, C.HIGH)
    C.digitalWrite (C.A2, C.LOW)
    C.digitalWrite (C.B1, C.HIGH)
    C.digitalWrite (C.B2, C.LOW)
    modifier := (float32(state.RTrigger) + 32768)/ 65536
    newPwm = int(float32(basePwm) * modifier)
    C.softPwmWrite(C.LPWM, C.int(float32(newPwm) * leftMod))
    C.softPwmWrite(C.RPWM, C.int(float32(newPwm) * rightMod))
  } else if state.LTrigger > -22767 || state.B {
    xbc.DEBUG("Backwards")
    C.digitalWrite (C.A1, C.LOW)
    C.digitalWrite (C.A2, C.HIGH)
    C.digitalWrite (C.B1, C.LOW)
    C.digitalWrite (C.B2, C.HIGH)
    modifier := (float32(state.LTrigger) + 32768)/ 65536
    newPwm := int(float32(basePwm) * modifier)
    C.softPwmWrite(C.LPWM, C.int(float32(newPwm) * leftMod))
    C.softPwmWrite(C.RPWM, C.int(float32(newPwm) * rightMod))
  } else if state.DPadX == 32767 {
    xbc.DEBUG("RIGHT")
    C.softPwmWrite(C.LPWM, 0);
  } else if state.DPadX == -32767 {
    xbc.DEBUG("LEFT")
    C.softPwmWrite(C.RPWM, 0);
  } else {
    xbc.DEBUG("STOP")
    C.softPwmWrite(C.LPWM, 0)
    C.softPwmWrite(C.RPWM, 0)
  }
  xbc.DEBUG("newPwm: ", newPwm);
  xbc.DEBUG("leftMod: ", leftMod);
  xbc.DEBUG("rightMod: ", rightMod);


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
    C.softPwmCreate(C.LPWM, 0, 100);
    C.softPwmCreate(C.RPWM, 0, 100);
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
