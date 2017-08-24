package main

import (
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
		if len(result) > 0 {
			metaData.Author = result[1]
			metaData.Title = result[2]
			metaData.PubYear, _ = strconv.Atoi(result[3])
			return metaData
		}
	}

	re, err = regexp.Compile(`(.+) - (.+)`)
	if err == nil {
		result := re.FindStringSubmatch(filename)
		if len(result) > 0 {
			metaData.Author = result[1]
			metaData.Title = result[2]
			return metaData
		}
	}

	re, err = regexp.Compile(`(.+) \((\d{4})\)`)
	if err == nil {
		result := re.FindStringSubmatch(filename)
		if len(result) > 0 {
			metaData.Title = result[1]
			metaData.PubYear, _ = strconv.Atoi(result[2])
			return metaData
		}
	}

	metaData.Title = filename

	return metaData
}
