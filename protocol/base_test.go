package protocol

import "testing"

func TestUint8(t *testing.T) {

	b := make([]byte, 0)
	WriteUint8(&b, 0xA9)
	if ReadUint8(&b) != 0xA9 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}

	b = make([]byte, 0)
	WriteUint8(&b, 0xA9)
	WriteUint8(&b, 0x11)
	WriteUint8(&b, 0x22)
	if ReadUint8(&b) != 0xA9 {
		t.Errorf("err")
	} else if len(b) != 2 {
		t.Errorf("err")
	}
	if ReadUint8(&b) != 0x11 {
		t.Errorf("err")
	} else if len(b) != 1 {
		t.Errorf("err")
	}
	if ReadUint8(&b) != 0x22 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}
	func() {
		defer func() {
			_ = recover()
		}()
		ReadUint8(&b)
		t.Errorf("err")
	}()

}

func TestUint32(t *testing.T) {

	b := make([]byte, 0)
	WriteUint32(&b, 0xA9Fa5168)
	if ReadUint32(&b) != 0xA9Fa5168 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}

	b = make([]byte, 0)
	WriteUint32(&b, 0xA9Fa5168)
	WriteUint32(&b, 0x11111111)
	WriteUint32(&b, 0x22222222)
	if ReadUint32(&b) != 0xA9Fa5168 {
		t.Errorf("err")
	} else if len(b) != 8 {
		t.Errorf("err")
	}
	if ReadUint32(&b) != 0x11111111 {
		t.Errorf("err")
	} else if len(b) != 4 {
		t.Errorf("err")
	}
	if ReadUint32(&b) != 0x22222222 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}
	func() {
		defer func() {
			_ = recover()
		}()
		ReadUint32(&b)
		t.Errorf("err")
	}()

}

func TestUint64(t *testing.T) {

	b := make([]byte, 0)
	WriteUint64(&b, 0xA9F1254215a5168)
	if ReadUint64(&b) != 0xA9F1254215a5168 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}

	b = make([]byte, 0)
	WriteUint64(&b, 0xA9F1254215a5168)
	WriteUint64(&b, 0x11111111)
	WriteUint64(&b, 0x22222222)
	if ReadUint64(&b) != 0xA9F1254215a5168 {
		t.Errorf("err")
	} else if len(b) != 16 {
		t.Errorf("err")
	}
	if ReadUint64(&b) != 0x11111111 {
		t.Errorf("err")
	} else if len(b) != 8 {
		t.Errorf("err")
	}
	if ReadUint64(&b) != 0x22222222 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}
	func() {
		defer func() {
			_ = recover()
		}()
		ReadUint64(&b)
		t.Errorf("err")
	}()

}
func TestInt32(t *testing.T) {

	b := make([]byte, 0)
	WriteInt32(&b, 0xA95168)
	if ReadInt32(&b) != 0xA95168 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}

	b = make([]byte, 0)
	WriteInt32(&b, -51676232)
	WriteInt32(&b, 0x11111111)
	WriteInt32(&b, 0x22222222)
	if ReadInt32(&b) != -51676232 {
		t.Errorf("err")
	} else if len(b) != 8 {
		t.Errorf("err")
	}
	if ReadInt32(&b) != 0x11111111 {
		t.Errorf("err")
	} else if len(b) != 4 {
		t.Errorf("err")
	}
	if ReadInt32(&b) != 0x22222222 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}
	func() {
		defer func() {
			_ = recover()
		}()
		ReadInt32(&b)
		t.Errorf("err")
	}()

}
func TestInt64(t *testing.T) {

	b := make([]byte, 0)
	WriteInt64(&b, 0x951681524147454)

	if len(b) != 8 {
		t.Error("len!=8")
	}

	if ReadInt64(&b) != 0x951681524147454 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}

	b = make([]byte, 0)
	WriteInt64(&b, -51674152456232)
	WriteInt64(&b, 0x11111111)
	WriteInt64(&b, 0x22222222)
	if ReadInt64(&b) != -51674152456232 {
		t.Errorf("err")
	} else if len(b) != 16 {
		t.Errorf("err")
	}
	if ReadInt64(&b) != 0x11111111 {
		t.Errorf("err")
	} else if len(b) != 8 {
		t.Errorf("err")
	}
	if ReadInt64(&b) != 0x22222222 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}
	func() {
		defer func() {
			_ = recover()
		}()
		ReadInt64(&b)
		t.Errorf("err")
	}()

}

func TestUint16(t *testing.T) {

	b := make([]byte, 0)
	WriteUint16(&b, 0xA9Fa)
	if ReadUint16(&b) != 0xA9Fa {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}

	b = make([]byte, 0)
	WriteUint16(&b, 0xA9Fa)
	WriteUint16(&b, 0x1111)
	WriteUint16(&b, 0x2222)
	if ReadUint16(&b) != 0xA9Fa {
		t.Errorf("err")
	} else if len(b) != 4 {
		t.Errorf("err")
	}
	if ReadUint16(&b) != 0x1111 {
		t.Errorf("err")
	} else if len(b) != 2 {
		t.Errorf("err")
	}
	if ReadUint16(&b) != 0x2222 {
		t.Errorf("err")
	} else if len(b) != 0 {
		t.Errorf("err")
	}
	func() {
		defer func() {
			_ = recover()
		}()
		ReadUint16(&b)
		t.Errorf("err")
	}()

}

func TestString(t *testing.T) {
	b := make([]byte, 0)
	defer func() {
		r := recover()
		if r != nil {
			t.Error(r)
		}
	}()

	Write2String(&b, "0123abc")
	s := Read2String(&b)

	if s != "0123abc" {
		t.Fatalf("")
	}
	if len(b) != 0 {
		t.Fatalf("")
	}

	Write2String(&b, "0123abc")
	Write2String(&b, "456def789")
	s = Read2String(&b)

	if s != "0123abc" {
		t.Fatalf("")
	}
	s = Read2String(&b)
	if s != "456def789" {
		t.Fatalf("")
	}
	if len(b) != 0 {
		t.Fatalf("")
	}

}
