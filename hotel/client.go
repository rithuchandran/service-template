package hotel

import (
	"compress/gzip"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type clientInt interface {
	getRegions() Regions
}

type client struct {
	url string
	*http.Client
}

func NewClient(url string) *client {
	return &client{url: url, Client: &http.Client{}}
}

func (client client) getRegions() Regions {
	request := createRequest(fmt.Sprintf("%s/regions", client.url))

	var regions Regions
	for ok := true; ok; {
		resp, err := client.Do(request)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("Do error", err)
			continue
		}
		request, ok = getNextLink(resp)

		decode(resp, &regions)

		_ = resp.Body.Close()
	}
	return regions
}

func decode(resp *http.Response, regions *Regions) {

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("gzip error", err)
		}
		resp.Body = gzipReader
	}
	err := json.NewDecoder(resp.Body).Decode(regions)
	if err != nil {
		fmt.Println("json decoding error", err)
	}
}

func getNextLink(resp *http.Response) (*http.Request, bool) {
	link := resp.Header.Get("Link")
	sep := func(c rune) bool {
		return c == ';' || c == '<' || c == '>' || c == '"'
	}
	if values := strings.FieldsFunc(link, sep); len(values) > 0 {
		return createRequest(values[0]), true
	}
	return &http.Request{}, false
}

func createRequest(target string) *http.Request {
	request, _ := http.NewRequest("GET", target, nil)

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Customer-Ip", "10.132.20.37") //to be taken from request
	request.Header.Add("User-Agent", "BigLife/0.1")
	request.Header.Add("Authorization", getAuthHeader())
	q := request.URL.Query()
	q.Add("language", "en-US")
	q.Add("include", "details")
	q.Add("include", "property_ids")
	q.Add("include", "property_ids_expanded")
	request.URL.RawQuery = q.Encode()

	return request
}

var now = func() time.Time {
	return time.Now()
}

func getAuthHeader() string {
	timeStamp := strconv.FormatInt(now().Unix(), 10)
	hash := sha512.New()
	apiKey := viper.GetString("API_KEY")
	hash.Write([]byte(apiKey + viper.GetString("SECRET_KEY") + timeStamp ))
	signature := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("EAN apikey=%s,signature=%s,timestamp=%s", apiKey, signature, timeStamp)
}
