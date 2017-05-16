package bixbar

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Bar is a single bar in system
type Bar interface {
	AddBlock(SimpleBlock)

	Start()

	Stop()
}

type bar struct {
	lock   sync.RWMutex
	blocks []SimpleBlock
	stop   bool
	tick   time.Duration
	reader io.Reader
	writer io.Writer

	iblocks map[string]InteractiveBlock

	close chan struct{}
}

func (b *bar) AddBlock(l SimpleBlock) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.blocks = append(b.blocks, l)
	if ib, ok := l.(InteractiveBlock); ok {
		b.iblocks[fmt.Sprintf("%s-%s", ib.Name(), ib.Instance())] = ib
	}
}

func (b *bar) writeLoop() {
	// First write the header
	v := header{
		Version:     1,
		ClickEvents: b.reader != nil,
		ContSignal:  "SIGCONT",
		StopSignal:  "SIGSTOP",
	}
	err := json.NewEncoder(b.writer).Encode(v)
	if err != nil {
		panic(err) // TODO : real error handling :)
	}
	var comma bool
theWriteLoop:
	for {
		select {
		case <-time.After(b.tick):
		case <-b.close:
			break theWriteLoop
		}
		b.write(comma)
		comma = true
	}
}

func (b *bar) callClick(c click) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if ib, ok := b.iblocks[fmt.Sprintf("%s-%s", c.Name, c.Instance)]; ok {
		ib.Click(c.X, c.Y, c.Button)
	}
}

func (b *bar) readLoop() {
	// TODO : write into channel to use select
	buf := []byte{' '}
	for {
		n, err := b.reader.Read(buf)
		if err != nil {
			break
		}

		if n < 1 {
			continue
		}
		if buf[0] != '[' && buf[0] != ',' {
			continue
		}
		c := click{}
		err = json.NewDecoder(b.reader).Decode(&c)
		if err != nil {
			break
		}
		b.callClick(c)
	}
}

func (b *bar) write(comma bool) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.stop {
		return
	}
	var all = make([]block, len(b.blocks))
	for i := range b.blocks {
		bl := block{
			FullText:            b.blocks[i].FullText(),
			ShortText:           b.blocks[i].ShortText(),
			Separator:           b.blocks[i].Separator(),
			Markup:              b.blocks[i].Markup(),
			Align:               b.blocks[i].Align(),
			MinWidth:            b.blocks[i].MinWidth(),
			SeparatorBlockWidth: b.blocks[i].SeparatorBlockWidth(),
			Urgent:              b.blocks[i].Urgent(),
		}
		if c, ok := b.blocks[i].Color(); ok {
			bl.Color = c.String()
		}
		if c, ok := b.blocks[i].Background(); ok {
			bl.Background = c.String()
		}
		if c, ok := b.blocks[i].Border(); ok {
			bl.Border = c.String()
		}

		if ib, ok := b.blocks[i].(InteractiveBlock); ok {
			bl.Name = ib.Name()
			bl.Instance = ib.Instance()
		}
		all[i] = bl
	}
	sep := []byte(",")
	if !comma {
		sep = []byte("[")
	}
	_, err := b.writer.Write(sep)
	if err != nil {
		panic(err) // TODO : :)
	}

	err = json.NewEncoder(b.writer).Encode(all)
	if err != nil {
		panic(err) // TODO : :)
	}
}

func (b *bar) sigWatch() {
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGCONT, syscall.SIGSTOP)

bigLoop:
	for {
		select {

		case s := <-sigs:
			b.lock.Lock()
			if s == syscall.SIGSTOP {
				b.stop = true
			}
			if s == syscall.SIGCONT {
				b.stop = false
			}
			b.lock.Unlock()
		case <-b.close:
			break bigLoop
		}
	}
}

func (b *bar) Start() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.close = make(chan struct{})

	go b.writeLoop()
	go b.sigWatch()
	if b.reader != nil {
		go b.readLoop()
	}
}

func (b *bar) Stop() {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.close == nil {
		return
	}

	close(b.close)
}

// NewBar create a bar and start it
func NewBar(d time.Duration, out io.Writer, in io.Reader) Bar {
	res := &bar{
		iblocks: make(map[string]InteractiveBlock),
		tick:    d,
		writer:  out,
	}
	if in != nil {
		res.reader = in
	}
	return res
}
