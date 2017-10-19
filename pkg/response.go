package tagc

type response struct {
	URL      string           `json:"url"`
	Meta     responseMeta     `json:"meta"`
	Elements responseElements `json:"elements,omitempty"`
}

func newResponse(url string, status int, contentType string, contentLength int, elements responseElements) response {
	return response{
		URL: url,
		Meta: responseMeta{
			Status:        status,
			ContentType:   contentType,
			ContentLength: contentLength,
		},
		Elements: elements,
	}
}

type responseMeta struct {
	Status        int    `json:"status"`
	ContentType   string `json:"content-type,omitempty"`
	ContentLength int    `json:"content-length"`
}

type responseElement struct {
	TagName string `json:"tag-name"`
	Count   int    `json:"count"`
}

type responseElements []responseElement

func newResponseElements(dup map[string]int) (re responseElements) {
	for tagName, cnt := range dup {
		re = append(re, responseElement{tagName, cnt})
	}
	return
}
