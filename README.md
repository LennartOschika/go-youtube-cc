# go-youtube-cc
For now only TTML works as a format. 
I started with implementing VVT, however the format is acutally trash unless you have hand made subtitles. 

To use this, just just call GetSubtitles with the youtube link that you want. The formats parameter does nothing.
You will get back a slice of the subtitles with the following fields:
timeStart: When does the caption occur in the video
timeEnd: Until what timee is the caption displayed
content: The actual content of the caption

I will (probably not) improve the code at some point to have it more robust, however currently it works fine for me. 

Stuff that could be improved:
- Better JSON and XML parsing
- Idk if it is possible to optimize to reduce the HTTP Request that is used to get the subtitle URL
- Implement the other formats (although for the videos that I tried I never saw any of the other 3 formats being used)
- YouTube could make it so that downloading captions doesn't require you to be the owner of the video
- Make more helper functions so that you can use my shit code in a less shitty way

Shoudout to ChatGPT for helping out with the XML and JSON parsing
