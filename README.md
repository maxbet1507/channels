# channels
channel functions for golang

## インストール

```
go get github.com/maxbet1507/channels
```

## 説明

よくある、無限キューチャンネルとPubSubです。

```go
func ExamplePubSub() {
	hub := PubSub()
	defer close(hub)

	hub <- "test1"

	sub1, closer1 := hub.Subscribe(true)
	defer closer1()

	sub2, closer2 := hub.Subscribe(true)
	defer closer2()

	hub <- "test2"

	fmt.Println(<-sub1)
	fmt.Println(<-sub2)

	// Output:
	// test2
	// test2
}
```