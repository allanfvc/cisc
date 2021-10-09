package github

import (
	"bytes"
	"math"
	"net/http"
)

func (g *GitHub) request(url, method string, body []byte) (*http.Response, error) {
	req, err := prepareRequest(url, method, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	return g.RestClient.Do(req)
}

func prepareRequest(url string, method string, body []byte) (*http.Request, error) {
	if body != nil {
		data := bytes.NewBuffer(body)
		return http.NewRequest(method, url, data)
	}
	return http.NewRequest(method, url, nil)
}

func getRuns(total, pagesize int) int {
	d := float64(total) / float64(pagesize)
	return int(math.Ceil(d))
}

func isZip(response *http.Response) bool {
	content := response.Header.Get("Content-Type")
	return content != "" && content == "application/zip"
}
