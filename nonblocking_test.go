package channels

import (
	"testing"
	"time"
)

func TestNonblocking(t *testing.T) {
	o, i := Nonblocking(false)

	i <- "test1"
	i <- "test2"

	if v, ok := <-o; !ok || v != "test1" {
		t.Fatal(v, ok)
	}
	if v, ok := <-o; !ok || v != "test2" {
		t.Fatal(v, ok)
	}

	i <- "test3"
	close(i)
	time.Sleep(100 * time.Millisecond)

	if v, ok := <-o; ok {
		t.Fatal(v, ok)
	}
}

func TestNonblockingMandatory(t *testing.T) {
	o, i := Nonblocking(true)

	i <- "test1"
	i <- "test2"

	if v, ok := <-o; !ok || v != "test1" {
		t.Fatal(v, ok)
	}
	if v, ok := <-o; !ok || v != "test2" {
		t.Fatal(v, ok)
	}

	i <- "test3"
	close(i)

	if v, ok := <-o; !ok || v != "test3" {
		t.Fatal(v, ok)
	}
	if v, ok := <-o; ok {
		t.Fatal(v, ok)
	}
}
