package main

import (
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "io"
  "net"
  "io/ioutil"
  "github.com/BurntSushi/toml"
  "flag"
  rc "github.com/jmittert/robotCar/lib"
  xbc "github.com/jmittert/xb360ctrl"
  "database/sql"
  _ "github.com/lib/pq"
)

type Config struct {
  ServerAddr string
  DbAddr string
  DbUser string
  DbPass string
  DbName string
}

var hw rc.HwState
var db *sql.DB

func cleanup(conn net.Conn, db *sql.DB, rc rc.HwState) {
  if conn != nil {
    conn.Close()
  }
  if db != nil {
    db.Close()
  }
  rc.A1 = 0
  rc.A2 = 0
  rc.B1 = 0
  rc.B2 = 0
  rc.LPWM = 0
  rc.RPWM = 0
  rc.Write()
}

func main() {
  checkFlags()
  config := readConfig()
  serverAddr := config.ServerAddr
  bytes := make([]byte, 6)

  var conn net.Conn
  // Catch SIGTERM and clean up properly
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  signal.Notify(c, syscall.SIGTERM)
  go func() {
    <-c
    cleanup(conn, db, hw)
    os.Exit(0)
  }()
  for {
    conn = connect(serverAddr)
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

func connectToDb(config Config) {
  user := config.DbUser
  addr := config.DbAddr
  name := config.DbName
  pass := config.DbName
  var err error
  db, err = sql.Open("postgresql",
  fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", user, pass, addr, name))
  checkError(err)
}
