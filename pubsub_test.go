package channels

import "testing"

func TestPubSub(t *testing.T) {
	hub := PubSub()

	hub <- "test1"

	sub1, closer1 := hub.Subscribe(true)
	defer closer1()

	sub2, closer2 := hub.Subscribe(true)
	defer closer2()

	hub <- "test2"

	if v, ok := <-sub1; !ok || v != "test2" {
		t.Fatal(v, ok)
	}
	if v, ok := <-sub2; !ok || v != "test2" {
		t.Fatal(v, ok)
	}

	closer2()
	hub <- "test3"

	if v, ok := <-sub1; !ok || v != "test3" {
		t.Fatal(v, ok)
	}
	if v, ok := <-sub2; ok {
		t.Fatal(v, ok)
	}

	hub <- "test4"
	close(hub)

	if v, ok := <-sub1; !ok || v != "test4" {
		t.Fatal(v, ok)
	}
	if v, ok := <-sub1; ok {
		t.Fatal(v, ok)
	}
}
