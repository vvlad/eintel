package main

import (
	"github.com/ogier/pflag"
	"github.com/vvlad/eintel"
)

func main() {

	pflag.Parse()
  eintel.PlayText("test")
  eintel.PlayText("test")
}
