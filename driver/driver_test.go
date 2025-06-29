package driver

import (
	"net"
	"testing"
	"time"
)

func TestDriverWorkerRegistration(t *testing.T) {
	d := New("localhost:0")
	if err := d.Start(); err != nil {
		t.Fatalf("failed to start driver: %v", err)
	}
	defer d.Stop()

	conn, err := net.Dial("tcp", d.Address())
	if err != nil {
		t.Fatalf("failed to dial driver: %v", err)
	}
	defer conn.Close()

	// give driver time to register
	time.Sleep(50 * time.Millisecond)

	if count := d.WorkerCount(); count != 1 {
		t.Fatalf("expected 1 worker, got %d", count)
	}
}
