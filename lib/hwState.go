package robotCar
/*
#cgo LDFLAGS: -lwiringPi -lpthread
#include <wiringPi.h>
#include <softPwm.h>
int A1 = 0;
int A2 = 1;
int LPWM = 5;
int B1 = 3;
int B2 = 4;
int RPWM = 6;
*/
import "C"

type GpioPin int
const (
  A1 GpioPin = iota
  A2
  B1
  B2
)

type PwmPin int
const (
  LPWM PwmPin = iota
  RPWM
)

var HIGH byte = 1
var LOW byte = 0

type HwState struct {
  A1 byte
  A2 byte
  B1 byte
  B2 byte
  LPWM uint8
  RPWM uint8
}

// Write writes the hardware state to the hardware
func (state *HwState) Write() {
  C.digitalWrite(C.A1, C.int(state.A1))
  C.digitalWrite(C.A2, C.int(state.A2))
  C.digitalWrite(C.B1, C.int(state.B1))
  C.digitalWrite(C.B2, C.int(state.B2))
  C.softPwmWrite(C.LPWM, C.int(state.LPWM))
  C.softPwmWrite(C.RPWM, C.int(state.RPWM))
}

func (state *HwState) Setup() {
    C.wiringPiSetup()
    C.pinMode(C.A1, C.OUTPUT)
    C.pinMode(C.A2, C.OUTPUT)
    C.pinMode(C.B1, C.OUTPUT)
    C.pinMode(C.B2, C.OUTPUT)
    C.softPwmCreate(C.LPWM, 0, 100);
    C.softPwmCreate(C.RPWM, 0, 100);
}

// MarshalBinary encodes the state into binary and returns the result
func (s *HwState) MarshalBinary() (data []byte, err error) {
  // 32 + 16 + 8 + 8 = 64bits = 8 bytes
  a := make([]byte, 6)

  a[0] = byte(s.A1)
  a[1] = byte(s.A2)
  a[2] = byte(s.B1)
  a[3] = byte(s.B2)
  a[4] = byte(s.LPWM)
  a[5] = byte(s.RPWM)
  return a, nil
}

// UnMarshalBinary unencodes a marshalled state
func (s *HwState) UnMarshalBinary(data []byte) (err error) {
  s.A1   = data[0]
  s.A2   = data[1]
  s.B1   = data[2]
  s.B2   = data[3]
  s.LPWM = data[4]
  s.RPWM = data[4]
  return nil
}
