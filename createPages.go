package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/slack-go/slack"
)

func createPages() error {
	paths, err := getLogfilePaths("logs")
	if err != nil {
		return err
	}

	channelList := []string{}
	for channel, channels := range paths {
		htemp := HistoryTemplate{channel, map[string]slack.History{}}
		for _, logfile := range channels {
			history, err := loadLogs(logfile)
			if err != nil {
				return err
			}
			htemp.History[logfile] = history
		}
		templating(channel, htemp)
		channelList = append(channelList, channel)
	}
	templatingIndexPage(channelList)
	return nil
}

func getLogfilePaths(dir string) (map[string][]string, error) {
	channels, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := map[string][]string{}
	for _, channel := range channels {
		if channel.IsDir() {
			paths[channel.Name()] = []string{}
			chandir := filepath.Join(dir, channel.Name())
			logs, err := ioutil.ReadDir(chandir)
			if err != nil {
				return nil, err
			}
			for _, log := range logs {
				if log.IsDir() {
					continue
				}
				paths[channel.Name()] = append(paths[channel.Name()], filepath.Join(chandir, log.Name()))
			}
			continue
		}
	}
	return paths, nil
}

func loadLogs(fname string) (slack.History, error) {
	var empty slack.History
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return empty, err
	}
	var data slack.History
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return empty, err
	}
	return data, nil
}

type HistoryTemplate struct {
	Channel string
	History map[string]slack.History
}

func templating(channel string, history HistoryTemplate) {
	var err error
	// channel page
	tpl := template.Must(template.ParseFiles("./template/channel_log.html"))
	// write
	fpath := "./pages"
	fname := channel + ".html"
	err = os.MkdirAll(fpath, 0755)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	file, err := os.Create(fpath + "/" + fname)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer file.Close()
	err = tpl.Execute(file, history)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func templatingIndexPage(channels []string) {
	var err error
	// index page
	tpl := template.Must(template.ParseFiles("./template/index.html"))
	file, err := os.Create("./index.html")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer file.Close()
	err = tpl.Execute(file, channels)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
