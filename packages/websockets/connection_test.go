package websockets_test

import (
	"testing"

	"github.com/bedrock/packages/websockets"
)

func TestConnectionManager_Add_Find(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewConnectionManager()
	conn := websockets.NewConn(nil, "app-1")

	mgr.Add(conn)

	found, ok := mgr.Find("app-1", conn.SocketID())

	if !ok {
		t.Fatal("expected Find to return true after Add")
	}

	if found.SocketID() != conn.SocketID() {
		t.Errorf("expected SocketID %q, got %q", conn.SocketID(), found.SocketID())
	}
}

func TestConnectionManager_Remove(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewConnectionManager()
	conn := websockets.NewConn(nil, "app-1")

	mgr.Add(conn)
	mgr.Remove("app-1", conn.SocketID())

	_, ok := mgr.Find("app-1", conn.SocketID())

	if ok {
		t.Error("expected Find to return false after Remove")
	}
}

func TestConnectionManager_Count(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewConnectionManager()
	c1 := websockets.NewConn(nil, "app-1")
	c2 := websockets.NewConn(nil, "app-1")
	c3 := websockets.NewConn(nil, "app-1")

	mgr.Add(c1)
	mgr.Add(c2)
	mgr.Add(c3)

	if count := mgr.Count("app-1"); count != 3 {
		t.Errorf("expected Count 3, got %d", count)
	}
}

func TestConnectionManager_All(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewConnectionManager()
	c1 := websockets.NewConn(nil, "app-1")
	c2 := websockets.NewConn(nil, "app-1")

	mgr.Add(c1)
	mgr.Add(c2)

	all := mgr.All("app-1")

	if len(all) != 2 {
		t.Errorf("expected 2 connections, got %d", len(all))
	}
}

func TestConnectionManager_MultipleApps(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewConnectionManager()
	connA := websockets.NewConn(nil, "app-a")
	connB := websockets.NewConn(nil, "app-b")

	mgr.Add(connA)
	mgr.Add(connB)

	if mgr.Count("app-a") != 1 {
		t.Errorf("expected 1 connection in app-a, got %d", mgr.Count("app-a"))
	}

	if mgr.Count("app-b") != 1 {
		t.Errorf("expected 1 connection in app-b, got %d", mgr.Count("app-b"))
	}

	// Confirm the connection in app-a cannot be found under app-b.
	_, ok := mgr.Find("app-b", connA.SocketID())

	if ok {
		t.Error("connection from app-a should not be found in app-b")
	}
}
