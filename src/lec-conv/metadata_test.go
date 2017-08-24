package main

import (
	"testing"
)

func testGetMetaData(t *testing.T, filename string, expected MetaData) {
	metaData := GetMetaData(filename)
	if metaData != expected {
		t.Errorf("actual: %s, expected: %s\n", metaData, expected)
	}

}

func TestGetMetaDataAuthorTitleYear(t *testing.T) {
	testGetMetaData(t, "A Foo - Foo Bar (2015)", MetaData{
		Author:  "A Foo",
		Title:   "Foo Bar",
		PubYear: 2015,
	})
}

func TestGetMetaDataAuthorTitle(t *testing.T) {
	testGetMetaData(t, "A-Foo - Foo Bar", MetaData{
		Author:  "A-Foo",
		Title:   "Foo Bar",
		PubYear: 0,
	})
}

func TestGetMetaDataTitleYear(t *testing.T) {
	testGetMetaData(t, "Foo-Bar (2015)", MetaData{
		Author:  "",
		Title:   "Foo-Bar",
		PubYear: 2015,
	})
}
