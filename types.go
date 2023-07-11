package main

import "encoding/xml"

type parsedSubtitle struct {
	timeStart string
	timeEnd   string
	content   string
}

// Shoutout ChatGPT
type TTML struct {
	XMLName xml.Name `xml:"tt"`
	Body    struct {
		Div struct {
			P []struct {
				Content string `xml:",chardata"`
				Begin   string `xml:"begin,attr"`
				End     string `xml:"end,attr"`
				Style   string `xml:"style,attr"`
			} `xml:"p"`
		} `xml:"div"`
	} `xml:"body"`
}
