package main

import (
	"fmt"
	"github.com/injoyai/bar"
	"time"
)

func main() {
	b := bar.New(
		bar.WithTotal(100),
		bar.WithFlush(),
		bar.WithFinal(func(b *bar.Bar) {
			fmt.Printf("\r\033[92m✓ 系统引导完毕！\033[0m\n")
		}),
	)
	for i := 0; i < 100; i++ {
		b.Add(1)
		b.Flush()
		time.Sleep(time.Millisecond * 100)
	}
}
