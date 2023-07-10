package main

import (
	"errors"
	"regexp"
	"strings"
)

func parseSubtitleString(subtitleString string, format Format) ([]parsedSubtitle, error) {
	var parsedSubtitle []parsedSubtitle
	var err error

	switch format {
	case FormatVTT:
		parsedSubtitle, err = parseSubtitleVTT(subtitleString)
	default:
		return nil, errors.New("Subtitle format not supported.")
	}
	if err != nil {
		return nil, err
	}

	return parsedSubtitle, nil
}

func parseSubtitleVTT(subtitleString string) ([]parsedSubtitle, error) {
	var returnSlice []parsedSubtitle
	var err error

	var sliceEntry parsedSubtitle

	subtitleSlice := strings.Split(subtitleString, "\n\n")

	headerLine, err := regexp.MatchString("^[^\\d]", subtitleSlice[0])
	if err != nil {
		return nil, err
	}
	if headerLine {
		subtitleSlice = append(subtitleSlice[1:])
	}

	for _, s := range subtitleSlice {
		if s == "" {
			continue
		}
		splitResult := strings.SplitN(s, "\n", 2)
		sliceEntry.timeStamp = splitResult[0]
		sliceEntry.content = strings.ReplaceAll(splitResult[1], "\n", " ")
		returnSlice = append(returnSlice, sliceEntry)
	}

	return returnSlice, err
}
