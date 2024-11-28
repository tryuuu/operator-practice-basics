package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func process(ctx context.Context) error {
	// 親コンテキストに3秒のタイムアウトを設定して子コンテキストを作成
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// 関数の終了時にキャンセルを呼び出し、リソースを解放
	defer cancel()

	sec := rand.Intn(6)
	fmt.Printf("wait %d sec: ", sec)

	// 処理の完了を通知するためのチャネルを作成(バッファサイズ1で非ブロッキング)
	done := make(chan error, 1)
	// gproutineで擬似プロセスを非同期に実行
	go func(sec int) {
		time.Sleep(time.Duration(sec) * time.Second)
		done <- nil
	}(sec)

	// select: 複数のチャネルを同時に監視し、どれか1つが準備できた場合にそのケースを実行
	select {
	case <-done: // タイムアウトしなかった場合
		fmt.Println("complete")
		return nil
	case <-ctx.Done(): // コンテキストがタイムアウトされた場合(チャネルが閉じられる)
		fmt.Println("timeout")
		return ctx.Err()
	}
}

func main() {
	// 乱数を生成するためのシードを設定
	rand.Seed(time.Now().UnixNano())
	// コンテキストを生成
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		err := process(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}
