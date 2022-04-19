package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/jkrus/Test_Simple_Multuplexor/pkg/models"
)

var reg = regexp.MustCompile(`(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w\.-]*)*\/?`)

func prepareRequest(r *http.Request) ([]string, error) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	urls := models.InfoDTO{}
	if err = json.Unmarshal(bytes, &urls); err != nil {
		return nil, err
	}
	for _, url := range urls.Urls {
		if !reg.Match([]byte(url)) {
			return []string{url}, fmt.Errorf("wrong url %v", url)
		}
	}

	if len(urls.Urls) > 20 {
		return nil, fmt.Errorf("too many URLs to request")
	}

	return urls.Urls, nil
}
