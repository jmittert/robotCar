package robotCar
import "fmt"
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

// write writes the given value to the specified pin. It returns true
// if the pin was updated, otherwise it returns false.
func (state *HwState) Write(pin GpioPin, value byte) (wrote bool, err error) {
  wrote = false
  switch pin {
  case A1:
    if state.A1 != value {
      C.digitalWrite (C.A1, C.int(value))
      wrote = true;
    }
  case A2:
    if state.A2 != value {
      C.digitalWrite (C.A2, C.int(value))
      wrote = true;
    }
  case B1:
    if state.B1 != value {
      C.digitalWrite (C.B1, C.int(value))
      wrote = true;
    }
  case B2:
    if state.B2 != value {
      C.digitalWrite (C.B2, C.int(value))
      wrote = true;
    }
  default:
    err = fmt.Errorf("%d is not a valid pin", pin)
  }
  return
}

// writePWM writes the given pwm value to the specified pin. It returns true
// if the pin was updated, otherwise it returns false.
func (state *HwState) WritePWM(pin PwmPin, value uint8) (wrote bool, err error) {
  wrote = false
  if value > 100 {
    return wrote, fmt.Errorf("%d is not a valid pwm value", value)
  }
  switch pin {
  case LPWM:
    if state.LPWM != value {
      C.softPwmWrite(C.LPWM, C.int(value))
    }
  case RPWM:
    if state.RPWM != value {
      C.softPwmWrite(C.RPWM, C.int(value))
    }
  default:
    err = fmt.Errorf("%d is not a valid pin", pin)
  }
  return
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
