package main

import (
	"github.com/aldernero/tui-apps/rule30/pkg/tui"
	"time"
)

func main() {
	tui.StartTea(time.Now().UnixNano())
}
