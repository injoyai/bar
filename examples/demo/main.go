package main

import (
	"github.com/injoyai/bar"
	"time"
)

func main() {
	x := bar.New(
		bar.WithTotal(60),
		bar.WithFormatDefault(func(p bar.Plan) {
			p.SetStyle("■")
		}),
	)
	for {
		time.Sleep(time.Millisecond * 100)
		x.Add(1)
		if x.Flush() {
			break
		}
	}
}
