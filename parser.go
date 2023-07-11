package subtitles

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func parseSubtitleResponse(htmlResponse *http.Response, format Format) ([]parsedSubtitle, error) {
	var parsedSubtitle []parsedSubtitle
	var err error

	switch format {
	case FormatVTT:
		parsedSubtitle, err = parseSubtitleVTT(htmlResponse)
	case FormatTTML:
		parsedSubtitle, err = parseSubtitleTTML(htmlResponse)
	default:
		return nil, errors.New("Subtitle format not supported.")
	}
	if err != nil {
		return nil, err
	}

	return parsedSubtitle, nil
}

func parseSubtitleVTT(httpResponse *http.Response) ([]parsedSubtitle, error) {
	var returnSlice []parsedSubtitle
	var err error

	defer httpResponse.Body.Close()

	b, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	subtitleString := string(b)

	var sliceEntry parsedSubtitle
	//V1: (?m)^(([\d:\.]+ --> [\d:\.]+)).*?$
	//V2: (?m)^(([\d:\.]+ --> [\d:\.]+)).*?\n
	cleanTimeStampRegex := regexp.MustCompile("(?m)^(([\\d:\\.]+ --> [\\d:\\.]+)).*?$")
	subtitleString = cleanTimeStampRegex.ReplaceAllString(subtitleString, "$1")

	cleanTextRegex := regexp.MustCompile("(?m)^(.*?</c><[\\d:\\.]+>.*?</c>)$")
	subtitleString = cleanTextRegex.ReplaceAllString(subtitleString, "")

	cleanDuplicatesRegex := regexp.MustCompile("(?m)^(([\\d:\\.]+ --> [\\d:\\.]+)\\n\\D.*)(\\s{4})")
	subtitleString = cleanDuplicatesRegex.ReplaceAllString(subtitleString, "")

	subtitleString = strings.SplitN(subtitleString, "\n \n\n\n", 2)[1]

	subtitleSlice := strings.Split(subtitleString, "\n\n\n")

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
		sliceEntry.timeStart = strings.Split(splitResult[0], " --> ")[0]
		sliceEntry.timeEnd = strings.Split(splitResult[0], " --> ")[1]
		sliceEntry.content = strings.ReplaceAll(splitResult[1], "\n", " ")
		returnSlice = append(returnSlice, sliceEntry)
	}

	return returnSlice, err
}

func parseSubtitleTTML(httpResponse *http.Response) ([]parsedSubtitle, error) {
	var returnSlice []parsedSubtitle
	var err error

	//defer httpResponse.Body.Close()

	b, err := io.ReadAll(httpResponse.Body)

	var subTitles TTML

	err = xml.Unmarshal(b, &subTitles)
	if err != nil {
		return nil, err
	}

	var sliceEntry parsedSubtitle
	for _, p := range subTitles.Body.Div.P {
		sliceEntry.timeStart = p.Begin
		sliceEntry.timeEnd = p.End
		sliceEntry.content = p.Content
		returnSlice = append(returnSlice, sliceEntry)
	}

	return returnSlice, err
}
