package main

import (
  "github.com/jmittert/xb360ctrl"
  "net"
  "os"
  "os/signal"
  "syscall"
  "time"
  "fmt"
)

func cleanup(fd int, conn net.Conn) {
  xb360ctrl.Close(fd)
  if conn != nil {
    conn.Close()
  }
}
func main() {
  xb360ctrl.DebugModeOn()
  service := ":2718"
  listener,_ := net.Listen("tcp", service)
  fd := xb360ctrl.Init("/dev/input/js0")
  var conn net.Conn

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  signal.Notify(c, syscall.SIGTERM)
  go func() {
    <-c
    cleanup(fd, conn)
    os.Exit(0)
  }()

  var loopCount float64 = 0
  t := time.Now()
  for {
    var err error
    conn, err = listener.Accept()
    if err != nil {
      continue
    }

    for {
      e := xb360ctrl.GetXbEvent(fd)
      bin, _ := e.MarshalBinary()
      //fmt.Println(bin)
      _, err = conn.Write(bin)
      if err != nil {
        break
      }
      loopCount++
      fmt.Printf("e/s: %f\r", loopCount / t.Sub(time.Now()).Seconds())
    }
  }
}
