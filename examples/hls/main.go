package main

import "github.com/injoyai/bar"

func main() {
	s := "http://devimages.apple.com.edgekey.net/streaming/examples/bipbop_4x3/gear2/prog_index.m3u8"
	bar.DownloadHLS(s, "./data")
}
