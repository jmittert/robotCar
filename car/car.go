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

type Config struct {
  ServerAddr string
}

var hwState HwState
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
        conn.Close()
        break;
      } else if count != 8 {
        conn.Close()
        break;
      }
      var e xbc.Xbc_event
      e.UnMarshalBinary(bytes)
      xbc.UpdateState(&e, &xbState)
      stateToHw(&xbState, &hwState)
    }
  }
}

func calcPWM(state *xbc.Xbc_state) (leftPwm uint8, rightPwm uint8){
  var basePwm float32 = 100
  var leftMod float32 = 1
  var rightMod float32 = 1
  if state.LStickX < 1000 {
    // <1000 -> Go right -> slow down right wheel
    leftMod -= float32(state.LStickX)/32768
    if leftMod < 0 {
      leftMod = 0
    }
  } else if state.LStickX > 1000 {
    // >1000 -> Go left -> slow down left wheel
    rightMod -= float32(state.LStickX)/-32768
    if rightMod < 0 {
      rightMod = 0
    }
  }
  if state.RTrigger > -22767 {
    modifier := (float32(state.RTrigger) + 32768)/ 65536
    basePwm = float32(basePwm) * modifier
  } else if state.LTrigger > -22767 {
    modifier := (float32(state.LTrigger) + 32768)/ 65536
    basePwm = float32(basePwm) * modifier
  }
  return uint8(leftMod*basePwm), uint8(rightMod*basePwm)
}

// Uses the current state of the controller to set the appropriate hw pins
func stateToHw(state *xbc.Xbc_state, hwState *HwState) {
  if state.RTrigger > -22767 {
    hwState.Write(A1, HIGH)
    hwState.Write(A2, LOW)
    hwState.Write(B1, HIGH)
    hwState.Write(B2, LOW)
  } else if state.LTrigger > -22767 {
    hwState.Write(A1, LOW)
    hwState.Write(A2, HIGH)
    hwState.Write(B1, LOW)
    hwState.Write(B2, HIGH)
  }
  lpwm, rpwm := calcPWM(state)
  hwState.WritePWM(RPWM, rpwm)
  hwState.WritePWM(LPWM, lpwm)
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
    hwState.Setup()
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
