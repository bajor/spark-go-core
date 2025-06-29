package driver

import (
	"encoding/json"
	"errors"
	"net"
	"sync"
)

type Driver struct {
	address  string
	listener net.Listener
	mu       sync.Mutex // guards listener and workers
	workers  []net.Conn
}

func New(address string) *Driver {
	return &Driver{address: address}
}

func (d *Driver) Start() error {
	ln, err := net.Listen("tcp", d.address)
	if err != nil {
		return err
	}

	d.mu.Lock()
	d.listener = ln
	d.workers = make([]net.Conn, 0)
	d.mu.Unlock()

	go d.acceptLoop()
	return nil
}

func (d *Driver) acceptLoop() {
	for {
		conn, err := d.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		d.mu.Lock()
		d.workers = append(d.workers, conn)
		d.mu.Unlock()
	}
}

func (d *Driver) Address() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.listener != nil {
		return d.listener.Addr().String()
	}
	return d.address
}

func (d *Driver) WorkerCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.workers)
}

func (d *Driver) Distribute(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	d.mu.Lock()
	workers := append([]net.Conn(nil), d.workers...)
	d.mu.Unlock()

	var firstErr error
	for _, w := range workers {
		if _, err := w.Write(data); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (d *Driver) Stop() {
	d.mu.Lock()
	ln := d.listener
	workers := append([]net.Conn(nil), d.workers...)
	d.listener = nil
	d.workers = nil
	d.mu.Unlock()

	if ln != nil {
		ln.Close()
	}
	for _, w := range workers {
		w.Close()
	}
}
