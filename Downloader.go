package subtitles

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
)

type Format string

const (
	FormatSBV  Format = "SBV"
	FormatSCC  Format = "SCC"
	FormatSRT  Format = "SRT"
	FormatTTML Format = "TTML"
	FormatVTT  Format = "VTT "
)

// Get video data
func DownloadVideo(link string) (*http.Response, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the subtitle from the response
func ExtractSubtitleURL(response *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return "", err
	}

	var ytPlayerResponse string

	doc.Find("body").Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		scriptContent := s.Text()

		regex := regexp.MustCompile(`^(?:var ytInitialPlayerResponse = )(.*);$`)
		match := regex.FindStringSubmatch(scriptContent)

		if len(match) > 1 {
			ytPlayerResponse = match[1]
			return false
		}
		return true
	})

	var ytPlayerResponseJSON map[string]interface{}

	json.Unmarshal([]byte(ytPlayerResponse), &ytPlayerResponseJSON)

	//There is probably some better way to do this
	baseUrl, ok := ytPlayerResponseJSON["captions"].(map[string]interface{})["playerCaptionsTracklistRenderer"].(map[string]interface{})["captionTracks"].([]interface{})[0].(map[string]interface{})["baseUrl"].(string)
	if !ok {
		return "", errors.New("Couldn't extract baseUrl from videoplayer data!")
	}

	return baseUrl, nil
}

// Get subtitles
// Does everything
func GetSubtitles(videoLink string, formats ...Format) ([]parsedSubtitle, error) {
	response, err := DownloadVideo(videoLink)
	if err != nil {
		return nil, err
	}

	baseURL, err := ExtractSubtitleURL(response)
	if err != nil {
		return nil, err
	}

	//TODO: implement different formats
	fullLink := baseURL + "&fmt=TTML"
	resp, err := http.Get(fullLink)
	if err != nil {
		return nil, err
	}

	parsedSubtitles, err := parseSubtitleResponse(resp, FormatTTML)

	return parsedSubtitles, nil
	//Build full URL
	//Download subtitle file(s)
	//Parse subtitle files
}

func main() {
	//Manual: https://www.youtube.com/watch?v=Q85l1Fenc5w
	//Auto generated: https://www.youtube.com/watch?v=6QQYUnSegaU&t=613s
	GetSubtitles("https://www.youtube.com/watch?v=6QQYUnSegaU&t=613s")

}
