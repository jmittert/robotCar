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

  if proc != nil {
    proc.Signal(syscall.SIGINT)
  }
}

type stateList struct {
  state rc.HwState
  next *stateList
}

func main() {
  checkFlags()
  config := readConfig()
  connectToDb(config)
  //startCamera()
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

  // Set up the database statements
  stateStmt, err := db.Prepare("INSERT INTO states (a1, a2, b1, b2, lpwm, rpwm) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT ON CONSTRAINT uniq DO UPDATE SET a1=states.a1 RETURNING id;")
  checkError(err)
  imgStmt, err := db.Prepare("INSERT INTO images (image, state1, state2, state3, state4, state5) VALUES($1, $2, $3, $4, $5, $6);")
  checkError(err)

  // Track the current and last 4 state ids, default to 0
  state := [5]int{0,0,0,0,0}
  currState := 0

  for {
    conn = connect(serverAddr)
    for {
      row := stateStmt.QueryRow(hw.A1, hw.A2, hw.B1, hw.B2, hw.LPWM, hw.RPWM)
      var id int
      err = row.Scan(&id)
      checkError(err)
      state[currState] = id

      var img []byte
      fileName := getLatestPic()
      if fileName != "" {
        img, err = ioutil.ReadFile(fileName)
        checkError(err)
        imgStmt.Exec(
          img,
          state[currState],
          state[(currState + 1) % 5],
          state[(currState + 2) % 5],
          state[(currState + 3) % 5],
          state[(currState + 4) % 5])
      }
      currState = (currState + 1) % 5
      count, err := conn.Read(bytes)
      // The new state we expect is 6 bytes
      if count != 6 || err == io.EOF {
        // On EOF, disconnect and look for another connection
        conn.Close()
        break
      }
      hw.UnMarshalBinary(bytes)
      hw.Write()
      // state
      conf := []byte{1}
      _, err = conn.Write(conf)
      if err != nil {
        fmt.Println(err)
        fmt.Println("!!!!!!")
        conn.Close()
        break
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
  attrs := new(os.ProcAttr)
  attrs.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
  os.Mkdir("/tmp/pics", os.ModeDir & os.ModeTemporary)
  args := []string{"-i", "/usr/local/lib/mjpg-streamer/input_raspicam.so -fps 2 -vf", "-o", "/usr/local/lib/mjpg-streamer/output_file.so -f /tmp/pics -s 5"}
  proc, err = os.StartProcess("/usr/local/bin/mjpg_streamer", args, attrs)
  if proc == nil {
    fmt.Println("proc is null!")
  } else {
    fmt.Println(*proc)
  }
  checkError(err)
}

func getLatestPic() string{
  files, err := ioutil.ReadDir("/tmp/pics")
  checkError(err)
  // We want the most recent file. Since the file are saved by
  // date, the last one will be the most recent
  numPics := len(files)
  if numPics == 0 {
    return ""
  }
  return "/tmp/pics/" + files[numPics - 1].Name()
}
