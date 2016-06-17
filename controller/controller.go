package main

import (
  "github.com/jmittert/xb360ctrl"
  "net"
)


func main() {
  service := ":2718"
  listener,_ := net.Listen("tcp", service)
  fd := xb360ctrl.Init("/dev/input/js0")
  for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }

    for {
      e := xb360ctrl.GetXbEvent(fd)
      bin, _ := e.MarshalBinary()
      conn.Write(bin)
    }
  }
  xb360ctrl.Close(fd)
}
