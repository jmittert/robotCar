package robotCar
import (
  "testing"
)

func TestMarshall(t *testing.T) {
  var e HwState
  e.A1 = 0;
  e.A2 = 0;
  e.B1 = 0;
  e.B2 = 0;
  e.LPWM = 0;
  e.RPWM = 0;

  data, _ := e.MarshalBinary()

  var e2 HwState
  e2.UnMarshalBinary(data)
  if e2.A1 != e.A1 {
    t.Error("Expected ", e.A1,", got ", e2.A1)
  }
  if e2.A2 != e.A2 {
    t.Error("Expected ", e.A2,", got ", e2.A2)
  }
  if e2.B1 != e.B1 {
    t.Error("Expected ", e.B1,", got ", e2.B1)
  }
  if e2.B2 != e.B2 {
    t.Error("Expected ", e.B2,", got ", e2.B2)
  }
  if e2.LPWM != e.LPWM {
    t.Error("Expected ", e.LPWM,", got ", e2.LPWM)
  }
  if e2.RPWM != e.RPWM {
    t.Error("Expected ", e.RPWM,", got ", e2.RPWM)
  }

  e.A1 = 1;
  e.A2 = 1;
  e.B1 = 1;
  e.B2 = 1;
  e.LPWM = 1;
  e.RPWM = 1;

  data, _ = e.MarshalBinary()

  e2.UnMarshalBinary(data)
  if e2.A1 != e.A1 {
    t.Error("Expected ", e.A1,", got ", e2.A1)
  }
  if e2.A2 != e.A2 {
    t.Error("Expected ", e.A2,", got ", e2.A2)
  }
  if e2.B1 != e.B1 {
    t.Error("Expected ", e.B1,", got ", e2.B1)
  }
  if e2.B2 != e.B2 {
    t.Error("Expected ", e.B2,", got ", e2.B2)
  }
  if e2.LPWM != e.LPWM {
    t.Error("Expected ", e.LPWM,", got ", e2.LPWM)
  }
  if e2.RPWM != e.RPWM {
    t.Error("Expected ", e.RPWM,", got ", e2.RPWM)
  }

  e.A1 = 1;
  e.A2 = 0;
  e.B1 = 1;
  e.B2 = 0;
  e.LPWM = 255;
  e.RPWM = 254;

  data, _ = e.MarshalBinary()

  e2.UnMarshalBinary(data)
  if e2.A1 != e.A1 {
    t.Error("Expected ", e.A1,", got ", e2.A1)
  }
  if e2.A2 != e.A2 {
    t.Error("Expected ", e.A2,", got ", e2.A2)
  }
  if e2.B1 != e.B1 {
    t.Error("Expected ", e.B1,", got ", e2.B1)
  }
  if e2.B2 != e.B2 {
    t.Error("Expected ", e.B2,", got ", e2.B2)
  }
  if e2.LPWM != e.LPWM {
    t.Error("Expected ", e.LPWM,", got ", e2.LPWM)
  }
  if e2.RPWM != e.RPWM {
    t.Error("Expected ", e.RPWM,", got ", e2.RPWM)
  }
}

