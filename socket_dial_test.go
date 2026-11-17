package vz

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func pipeConnection(t *testing.T) (*VirtioSocketConnection, net.Conn) {
	t.Helper()
	ours, theirs := net.Pipe()
	t.Cleanup(func() { _ = ours.Close(); _ = theirs.Close() })
	return &VirtioSocketConnection{rawConn: ours}, theirs
}

// closed reports whether the connection whose other end is peer was closed.
func closed(peer net.Conn) bool {
	_ = peer.SetReadDeadline(time.Now().Add(time.Second))
	_, err := peer.Read(make([]byte, 1))
	return errors.Is(err, io.EOF)
}

func TestPendingDialReturnsTheCompletion(t *testing.T) {
	conn, _ := pipeConnection(t)
	dial := newPendingDial()
	go dial.complete(conn, nil)
	got, err := dial.wait(context.Background())
	if err != nil || got != conn {
		t.Fatalf("wait = %v, %v; want the completed connection", got, err)
	}

	failed := errors.New("no listener")
	dial = newPendingDial()
	go dial.complete(nil, failed)
	if got, err := dial.wait(context.Background()); got != nil || !errors.Is(err, failed) {
		t.Fatalf("wait = %v, %v; want the completion's error", got, err)
	}
}

func TestPendingDialReturnsWhenItsContextEnds(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	dial := newPendingDial()
	start := time.Now()
	conn, err := dial.wait(ctx)
	if conn != nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("wait = %v, %v; want the context's error", conn, err)
	}
	if took := time.Since(start); took > time.Second {
		t.Fatalf("wait outlived its context by %s", took)
	}
}

func TestPendingDialClosesAConnectionThatCompletesLate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	dial := newPendingDial()
	if _, err := dial.wait(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("wait = %v; want context.Canceled", err)
	}

	conn, peer := pipeConnection(t)
	done := make(chan struct{})
	go func() { dial.complete(conn, nil); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("the completion handler blocked on a caller that has gone")
	}
	if !closed(peer) {
		t.Fatal("a connection that completed after its caller left was not closed")
	}
}

func TestPendingDialClosesACompletionItRacedWith(t *testing.T) {
	// The handler delivers, then the caller's context ends before it reads:
	// whichever way wait goes, the connection is returned or closed.
	conn, peer := pipeConnection(t)
	dial := newPendingDial()
	dial.complete(conn, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, err := dial.wait(ctx)
	switch {
	case err == nil && got == conn:
	case errors.Is(err, context.Canceled) && closed(peer):
	default:
		t.Fatalf("wait = %v, %v, and the connection was left open", got, err)
	}
}
