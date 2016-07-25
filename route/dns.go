package route

import (
	"fmt"
	"net"

	"golang.org/x/text/encoding/simplifiedchinese"

	"os/exec"

	"bufio"
	"bytes"

	"regexp"

	"os"
	"path/filepath"

	"encoding/gob"

	"strings"

	"syscall"

	"strconv"

	log "github.com/gamexg/log4go"
)

var RE_IName, _ = regexp.Compile(`".+?"`)
var RE_DHCP, _ = regexp.Compile(`DHCP`)

var RE_IIp, _ = regexp.Compile(`\d+\.\d+\.\d+\.\d+`)

func SetDns(tunName string, dnsIps []net.IP) error {
	dnss, err := getDns()
	if err != nil {
		return err
	}

	if len(dnss) == 0 {
		return fmt.Errorf("未找到网卡。")
	}

	bDna := make([]IDns, 0)
	var vpnDns *IDns

	for i := 0; i < len(dnss); i++ {
		d := &dnss[i]
		log.Debug("找到网卡 %v ，vpn网卡名字%v\r\n", d.Name, tunName)
		if d.Name == tunName {
			vpnDns = d
		} else if len(d.Ips) != 0 {
			bDna = append(bDna, *d)
		}
	}

	if vpnDns == nil {
		vpnDns = &IDns{
			Name: tunName,
			Dhcp: true,
		}
	}

	if err := bakDns(append(bDna, *vpnDns)); err != nil {
		return fmt.Errorf("备份 dns 记录失败,%v", err)
	}

	log.Debug("修改网卡 %v 的dns为 %#v ", vpnDns.Name, dnsIps)

	if err := setDns(vpnDns.Name, false, dnsIps); err != nil {
		return err
	}

	for _, d := range bDna {
		if err := setDns(d.Name, false, []net.IP{net.IPv4(0, 0, 0, 0)}); err != nil {
			log.Error("删除网卡 %v 的dns记录失败，%v", d.Name, err)
		}
	}

	return nil
}

func getBakDnsFilePath() string {
	myPath, _ := exec.LookPath(os.Args[0])
	absPath, _ := filepath.Abs(myPath)
	dir := filepath.Dir(absPath)

	return filepath.Join(dir, "dns.bak")
}

func bakDns(d []IDns) error {
	bakPath := getBakDnsFilePath()
	f, err := os.Create(bakPath)
	if err != nil {
		return err
	}
	defer f.Close()

	en := gob.NewEncoder(f)

	err = en.Encode(d)

	if err != nil {
		return err
	}

	return nil
}

func ResetDns() error {
	bakPath := getBakDnsFilePath()
	if Exist(bakPath) {
		log.Debug("找到 dns 备份文件 %v ，开始还原。", bakPath)
		f, err := os.Open(bakPath)
		if err != nil {
			return err
		}
		defer os.Remove(bakPath)
		defer f.Close()

		de := gob.NewDecoder(f)

		var ds []IDns
		if err := de.Decode(&ds); err != nil {
			return err
		}

		for _, d := range ds {
			if d.Name != "" {
				if err := setDns(d.Name, d.Dhcp, d.Ips); err != nil {
					log.Error("%v", err)
				}
			}
		}
	} else {
		log.Debug("不存在 dns 备份文件 %v ，无需还原。", bakPath)
	}
	return nil
}

func setDns(name string, dhcp bool, ips []net.IP) error {
	log.Debug("设置网卡 %v 的dns为dhcp(%v) %#v", name, dhcp, ips)

	if dhcp {
		return cmdRun("netsh", "interface", "ip", "set", "dns", fmt.Sprintf(`name=%v`, name), `source=dhcp`)
	}

	if err := cmdRun("netsh", "interface", "ip", "set", "dns", fmt.Sprintf(`name=%v`, name), `source=static`, `addr=none`); err != nil {
		log.Error("清空 tap dns 错误，%v", err)
	}

	for i := 0; i < len(ips); i++ {
		if err := cmdRun(`netsh`, `interface`, `ip`, `add`, `dns`, fmt.Sprintf(`name=%v`, name), fmt.Sprintf(`addr=%v`, ips[i].String()), `index=`+strconv.Itoa(i+1)); err != nil {
			log.Error("设置 tap dns 错误，%v", err)
		}
	}
	return nil
}

type IDns struct {
	Name string
	Ips  []net.IP
	Dhcp bool
}

func getDns() ([]IDns, error) {
	res := make([]IDns, 0)

	cmd := exec.Command("netsh", "interface", "ip", "show", "dns")
	o, err := cmd.Output()
	if err == nil {
		out, err2 := simplifiedchinese.GB18030.NewDecoder().Bytes(o)
		if err2 == nil {
			o = out
		}

		d := IDns{}

		save := func() {
			if d.Name != "" {
				res = append(res, d)
				d = IDns{}
			}
		}

		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			text := scanner.Text()
			// 如果包含 ""则检查 d ，保存并新建
			if name := RE_IName.FindString(text); name != "" {
				save()
				d.Name = strings.Trim(name, `"`)
			}

			// 如果包含ip就添加到 d
			if ip := RE_IIp.FindString(text); ip != "" {
				_ip := net.ParseIP(ip)
				if _ip != nil {
					d.Ips = append(d.Ips, _ip)
				}
			}

			// 检查是否是dhcp
			if dhcp := RE_DHCP.FindString(text); dhcp != "" {
				d.Dhcp = true
			}
		}
		//最后在检查下d，添加到res
		save()

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		log.Debug("找到如下 dns 记录：%#v", res)

		return res, nil
	} else {
		return nil, fmt.Errorf("获得dns错误，%v ,%v", o, err)
	}
}

func cmdRun(cmd string, arg ...string) error {
	setip := exec.Command(cmd, arg...)
	setip.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}

	o, err := setip.Output()
	if o != nil {
		out, err2 := simplifiedchinese.GB18030.NewDecoder().Bytes(o)
		if err2 == nil {
			o = out
		}
	}
	if err != nil {
		return fmt.Errorf("cmd:%s,arg:%v,err:%s,out:%s", cmd, arg, err, string(o))
	}
	log.Debug("cmd:%s,arg:%v,err:%s,out:%s", cmd, arg, err, string(o))
	return nil
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
