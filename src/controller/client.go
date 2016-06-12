package main

import (
  "fmt"
  "os"
  "net"
)

func main() {
  service := "127.0.0.1:2718"

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

    for {
      char := make([]byte, 1)
      _, err = conn.Read(char)
      print(string(char))
      if err != nil {
        fmt.Println("\nConnection lost")
        break
      }
    }
  }
  os.Exit(0)
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}

