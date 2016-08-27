package main

import (
  "net"
  "os"
  "os/signal"
  "syscall"
  rc "github.com/jmittert/robotCar/lib"
  xbc "github.com/jmittert/xb360ctrl"
  "time"
  "fmt"
)

func cleanup(fd int, conn net.Conn) {
  xbc.Close(fd)
  if conn != nil {
    conn.Close()
  }
}

func main() {
  service := ":2718"
  listener,_ := net.Listen("tcp", service)
  fd := xbc.Init("/dev/input/js0")
  var conn net.Conn

  // Catch SIGTERM and clean up properly
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  signal.Notify(c, syscall.SIGTERM)
  go func() {
    <-c
    cleanup(fd, conn)
    os.Exit(0)
  }()

  for {
    var err error
    conn, err = listener.Accept()
    if err != nil {
      continue
    }
    fmt.Println("Connected")

    start := time.Now()
    var count int64 = 0
    for {
      // Update the current state
      st := getState(fd)

      // Send the state
      _, err = conn.Write(st)
      if err != nil {
        break
      }

      // Wait for confirmation before continuing
      var confirm []byte
      fmt.Println("_")
      _, err = conn.Read(confirm)
      for err != nil {
        _, err = conn.Read(confirm)
        //fmt.Println("?")
      }
      fmt.Println("!")
      count++
      fmt.Println(time.Now().Sub(start).Nanoseconds()/count/1000000)
    }
  }
}

func getState(fd int) []byte {
  var xbState xbc.Xbc_state
  xbc.PrepState(&xbState)
  var hwState rc.HwState
  e := xbc.GetXbEvent(fd)
  xbc.UpdateState(e, &xbState)
  stateToHw(&xbState, &hwState)
  bin, _ := hwState.MarshalBinary()
  return bin
}

func calcPWM(state *xbc.Xbc_state) (leftPwm uint8, rightPwm uint8){
  var basePwm float32 = 100
  var leftMod float32 = 1
  var rightMod float32 = 1
  if state.LStickX > 1000 {
    // <1000 -> Go right -> slow down right wheel
    leftMod -= float32(state.LStickX)/32768
    if leftMod < 0 {
      leftMod = 0
    }
  } else if state.LStickX < -1000 {
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
  } else {
    basePwm = 0
  }
  return uint8(leftMod*basePwm), uint8(rightMod*basePwm)
}

// Uses the current state of the controller to set the appropriate hw pins
func stateToHw(state *xbc.Xbc_state, hwState *rc.HwState) {
  if state.RTrigger > -22767 {
    hwState.A1 = rc.HIGH
    hwState.A2 = rc.LOW
    hwState.B1 = rc.HIGH
    hwState.B2 = rc.LOW
  } else if state.LTrigger > -22767 {
    hwState.A1 = rc.LOW
    hwState.A2 = rc.HIGH
    hwState.B1 = rc.LOW
    hwState.B2 = rc.HIGH
  }
  lpwm, rpwm := calcPWM(state)
  hwState.LPWM = lpwm
  hwState.RPWM = rpwm
}
