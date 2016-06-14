package main

import ("fmt"
"github.com/jmittert/joystick"
"net"
"strconv"
)

func main() {
  service := ":2718"
  listener,_ := net.Listen("tcp", service)
  joystick.Init();
  for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }
    fmt.Println("Got Connection")

    for {
      event := joystick.GetJsEvent()
      fmt.Println(strconv.Itoa(event.Number), " ", strconv.Itoa(event.Value))
      conn.Write([]byte(strconv.Itoa(event.Number)))
      conn.Write([]byte(" "))
      conn.Write([]byte(strconv.Itoa(event.Value)))
      conn.Write([]byte("\n"))
    }
  }
}
