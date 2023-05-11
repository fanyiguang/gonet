package route

//func TestSetDns(t *testing.T) {
//	if err := SetDns([]net.IP{net.IPv4(11, 12, 13, 14)}); err != nil {
//		t.Error(err)
//	}
//
//	// 测试dns
//	cmd := exec.Command("NSLookup", "localhost")
//	o, err := cmd.Output()
//	if err != nil {
//		t.Error(err)
//	}
//	if o != nil {
//		if strings.Contains(string(o), "11.12.13.14") == false {
//			t.Error("设置dns服务器无效：%v", string(o))
//		}
//	} else {
//		t.Error("o==nil")
//	}
//
//	if err := ResetDns(); err != nil {
//		t.Error(err)
//	}
//	cmd = exec.Command("NSLookup", "localhost")
//	o, err = cmd.Output()
//	if err != nil {
//		t.Error(err)
//	}
//	if o != nil {
//		if strings.Contains(string(o), "11.12.13.14") == true {
//			t.Error("复位无效：%v", string(o))
//		}
//	} else {
//		t.Error("o==nil")
//	}
//
//}
