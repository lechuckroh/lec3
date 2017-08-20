package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type MetaData struct {
	Title   string
	Author  string
	PubYear int
}

func GetMetaData(filename string) MetaData {
	metaData := MetaData{}

	re, err := regexp.Compile(`(.+) - (.+) \((\d{4})\)`)
	if err == nil {
		result := re.FindStringSubmatch(filename)
		metaData.Author = result[1]
		metaData.Title = result[2]
		metaData.PubYear, _ = strconv.Atoi(result[3])
	} else {
		fmt.Printf(err.Error())
	}

	return metaData
}
