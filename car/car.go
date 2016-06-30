package main

import (
  "fmt"
  "os"
  "io"
  "net"
  "io/ioutil"
  "github.com/BurntSushi/toml"
  "flag"
  rc "github.com/jmittert/robotCar/lib"
  xbc "github.com/jmittert/xb360ctrl"
)

type Config struct {
  ServerAddr string
}

var hw rc.HwState
func main() {
  checkFlags()
  config := readConfig()
  serverAddr := config.ServerAddr
  bytes := make([]byte, 6)

  for {
    conn := connect(serverAddr)
    for {
      count, err := conn.Read(bytes)
      if count != 6 || err == io.EOF {
        // On EOF, disconnect and look for another connection
        conn.Close()
        break;
      }
      hw.UnMarshalBinary(bytes)
      hw.Write()
    }
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
    hw.Setup()
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
