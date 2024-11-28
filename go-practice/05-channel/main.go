package main

import (
	"fmt"
	"time"
)

func process(ch chan int) {
	fmt.Println("process start")
	time.Sleep(1 * time.Second)
	ch <- 1 // チャネルに送信
	fmt.Println("process end")
}

func main() {
	// mainが先に終わってしまわないように、チャネルを使って待機する
	waitCh := make(chan int) // チャネルを初期化

	go process(waitCh)

	fmt.Println("waiting")

	<-waitCh // チャネルから受信されるまで待機

	fmt.Println("done")
}
