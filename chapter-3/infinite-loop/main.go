package main

import "fmt"

// GOMAXPROCS=1 GODEBUG=asyncpreemptoff=1 を使って呼び出すとinfiniteの関数は呼ばれない

func main() {
	go infinite()	
	var i int
	for {
		// システムコール呼び出しのタイミングでgoroutineのスイッチが行われる
		// fmt.Println("main")
		// 以下のようなシステうコールを伴わないような関数呼び出しの場合は相変わらずgoroutineのスイッチが行われないため、infiniteの処理が行われない
		i = increment(i)
	}
}

func increment(i int) int {
	return i+1
}

func infinite() {
	for {
		fmt.Println("infinite!")
	}
}