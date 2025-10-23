package main

import (
	"github.com/injoyai/bar"
	"time"
)

func main() {
	b := bar.New(
		bar.WithTotal(100),
		bar.WithFormat(
			bar.WithAnimationMoon(),
			bar.WithRate(),
		),
	)
	for i := 0; i < 100; i++ {
		b.Add(1)
		b.Flush()
		time.Sleep(time.Millisecond * 100)
	}
}
