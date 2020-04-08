package main

import (
	"fmt"
	"strconv"
	"time"
)

func unix2str(n string, m string) string {
	nsec, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		panic(err)
	}
	msec, err := strconv.ParseInt(m, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(nsec, msec)
	return tm.Format("15:04:05")
}

type MessageTemplate struct {
	Time    string
	Name    string
	Message string
}

type Template struct {
	Channel  string
	Messages []MessageTemplate
}

func main() {
	var err error
	err = saveLog()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = createPages()
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
