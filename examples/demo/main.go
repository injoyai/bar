package main

import (
	"time"

	"github.com/injoyai/bar"
)

func main() {
	x := bar.New(
		bar.WithTotal(60),
		bar.WithFormatDefault(func(p *bar.Plan) {
			p.SetStyle("■")
			p.SetPadding(".")
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
