package main

import (
	"time"

	"github.com/injoyai/bar"
)

func main() {
	x := bar.New(
		bar.WithTotal(60000),
		bar.WithFormatDefault(func(p *bar.Plan) {
			p.SetStyle("■")
			p.SetPadding(".")
		}),
	)

	for range 100 {
		go func() {
			for {
				time.Sleep(time.Millisecond * 10)
				if x.Add(1).Flush().Closed() {
					break
				}
			}
		}()
	}

	<-x.Done()

}
