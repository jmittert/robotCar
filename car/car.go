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
  "time"
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
var conn net.Conn
var proc *os.Process

func cleanup() {
  if conn != nil {
    conn.Close()
  }
  if db != nil {
    db.Close()
  }
  hw.A1 = 0
  hw.A2 = 0
  hw.B1 = 0
  hw.B2 = 0
  hw.LPWM = 0
  hw.RPWM = 0
  hw.Write()

  proc.Signal(syscall.SIGINT)
}

func main() {
  checkFlags()
  config := readConfig()
  connectToDb(config)
  serverAddr := config.ServerAddr
  bytes := make([]byte, 6)

  var conn net.Conn
  // Catch SIGTERM and clean up properly
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  signal.Notify(c, syscall.SIGTERM)
  go func() {
    <-c
    cleanup()
    os.Exit(0)
  }()
  stmt, err := db.Prepare("INSERT INTO images (image, a1, a2, b1, b2, lpwm, rpwm) VALUES($1, $2, $3, $4, $5, $6, $7);")
  checkError(err)
  for {
    conn = connect(serverAddr)
    last := time.Now()
    for {
      count, err := conn.Read(bytes)
      if count != 6 || err == io.EOF {
        // On EOF, disconnect and look for another connection
        conn.Close()
        break;
      }
      hw.UnMarshalBinary(bytes)
      hw.Write()
      // Save the state every half second
      if time.Since(last).Nanoseconds() > 500000000 {
        var img []byte
        img, err = ioutil.ReadFile(getLatestPic())
        checkError(err)
        stmt.Exec(img, hw.A1, hw.A2, hw.B1, hw.B2, hw.LPWM, hw.RPWM)
        last = time.Now()
      }
    }
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    cleanup();
    os.Exit(1)
  }
}

// checkFlags is a help function to read and handle the command line flags
func checkFlags() {
  dbgPtr := flag.Bool("dbg", false, "Set debug mode")
  hwPtr  := flag.Bool("hw", true, "Enable talking with hardware")

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
  f, err := ioutil.ReadFile("/etc/carrc");
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
  db, err = sql.Open("postgres",
  fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, addr, name))
  checkError(err)
}

func startCamera() {
  var err error
  args := []string{"-i", "input_raspicam.so -fps 2 -vf", "-o", "output_file.so -f /tmp/pics -s 5"}
  proc, err = os.StartProcess("mjpg_streamer", args, nil)
  checkError(err)
}

func getLatestPic() string{
  files, err := ioutil.ReadDir("/tmp/pics")
  checkError(err)
  // We want the most recent file. Since the file are saved by
  // date, the last one will be the most recent
  return files[len(files)].Name()
}
