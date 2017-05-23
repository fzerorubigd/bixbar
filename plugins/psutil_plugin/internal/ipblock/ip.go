package ipblock

import (
	"html/template"
	"sync"

	"net"

	"bytes"

	"github.com/fzerorubigd/bixbar"
	psnet "github.com/shirou/gopsutil/net"
)

const ipDefaultTemplate = `{{ .Name }}:{{ .Addr }}`

type ipTemplateData struct {
	Name         string
	Addr         string
	MTU          int
	HardwareAddr string
	UP           bool
	LoopBack     bool
	MultiCast    bool
	P2P          bool
}

type ipData struct {
	tpl     *template.Template
	iface   string
	data    psnet.InterfaceStat
	addrIdx int
	err     error
	text    string

	lock sync.RWMutex
}

func (ip *ipData) render() {
	ip.lock.Lock()
	defer ip.lock.Unlock()

	buf := &bytes.Buffer{}

	contain := func(s string) bool {
		for _, i := range ip.data.Flags {
			if i == s {
				return true
			}
		}
		return false
	}

	err := ip.tpl.Execute(buf, ipTemplateData{
		Name:         ip.data.Name,
		Addr:         ip.data.Addrs[ip.addrIdx].Addr,
		MTU:          ip.data.MTU,
		HardwareAddr: ip.data.HardwareAddr,
		UP:           contain("up"),
		LoopBack:     contain("loopback"),
		MultiCast:    contain("multicast"),
		P2P:          contain("pointtopoint"),
	})
	if err != nil {
		ip.err = err
	}

	ip.text = buf.String()
}

func (ip *ipData) getText() string {
	ip.lock.RLock()
	defer ip.lock.RUnlock()

	if ip.err != nil {
		return ip.err.Error()
	}

	return ip.text
}

func (ip *ipData) FullText() string {
	return ip.getText()
}

func (ip *ipData) ShortText() string {
	return ip.getText()
}

func (ip *ipData) MinWidth() bixbar.StringInt {
	return bixbar.StringInt{String: ip.getText()}
}

func (ip *ipData) Align() bixbar.Align {
	return bixbar.Align("left")
}

func (ip *ipData) Color() (*bixbar.Color, bool) {
	ip.lock.RLock()
	defer ip.lock.RUnlock()

	if ip.err != nil {
		red, _ := bixbar.NewColor("#FF0000")
		return red, true
	}
	return nil, false
}

func (ip *ipData) Background() (*bixbar.Color, bool) {
	return nil, false
}

func (ip *ipData) Border() (*bixbar.Color, bool) {
	return nil, false
}

func (ip *ipData) Separator() bool {
	return true
}

func (ip *ipData) SeparatorBlockWidth() int {
	return 15
}

func (ip *ipData) Urgent() bool {
	ip.lock.RLock()
	defer ip.lock.RUnlock()

	return ip.err != nil
}

func (ip *ipData) Markup() bixbar.Markup {
	return bixbar.Markup("none")
}

func (ip *ipData) Click(int, int, bixbar.Button) {
	ip.lock.Lock()
	ip.addrIdx++
	if ip.addrIdx >= len(ip.data.Addrs) {
		ip.addrIdx = 0
	}
	ip.lock.Unlock()
	ip.render()
}

// NewIPBlock return a new ip block
func NewIPBlock(iface string, tplString string) bixbar.InteractiveBlock {
	if iface == "" {
		ifi := RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
		if ifi != nil {
			iface = ifi.Name
		}
	}

	if tplString == "" {
		tplString = ipDefaultTemplate
	}

	idata := ipData{
		iface: iface,
	}

	all, err := psnet.Interfaces()
	if err == nil {
		idata.tpl, err = template.New("ipdata").Parse(tplString)
	}

	if err != nil {
		idata.err = err
		return &idata
	}
	idata.err = err
	for i := range all {
		if all[i].Name == idata.iface {
			idata.data = all[i]
			idata.addrIdx = 0
			break
		}
	}
	idata.render()
	return &idata
}
