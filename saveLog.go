package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

func getToday() time.Time {
	ts := time.Now()
	ts = ts.Add(time.Hour * -24)
	return ts.Round(time.Hour * 24)
}

func saveLog() error {
	api := slack.New(os.Getenv("SLACK_API_KEY"))
	channels, err := api.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	historyparam := slack.NewHistoryParameters()
	today := getToday()
	yesterday := today.Add(time.Hour * -24)
	historyparam.Oldest = strconv.FormatInt(yesterday.Unix(), 10)
	historyparam.Latest = strconv.FormatInt(today.Unix(), 10)

	fmt.Printf("time : %s -> %s\n", yesterday, today)
	for _, channel := range channels {
		fmt.Printf("ID: %s, Name: %s\n", channel.ID, channel.Name)
		history, err := api.GetChannelHistory(channel.ID, historyparam)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		bjson, err := json.Marshal(history)

		// write
		fpath := "./logs/" + channel.Name
		fname := today.Format("20060102") + ".html"
		err = os.MkdirAll(fpath, 0755)
		if err != nil {
			fmt.Printf("%s\n", err)
			return err
		}
		file, err := os.Create(fpath + "/" + fname)
		if err != nil {
			fmt.Printf("%s\n", err)
			return err
		}
		defer file.Close()
		_, err = file.Write(bjson)
		if err != nil {
			fmt.Printf("%s\n", err)
			return err
		}
	}
	return nil
}
