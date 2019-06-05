package hotel

import (
	"compress/gzip"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

func TestGetRegions(t *testing.T) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", "regions_stub.txt"))
	if err != nil {
		panic(err)
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write(b)
	})
	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := NewClient("http://test_url")
	client.Client = httpCli

	regions := client.getRegions()

	assert.Equal(t, 250, len(regions))
	assert.Equal(t, "Albania", regions["2"].Name)
}

func TestAuthorization(t *testing.T) {
	viper.Set("API_KEY", "abc")
	viper.Set("SECRET_KEY", "secret")
	now = func() time.Time {
		return time.Unix(1559215747, 0)
	}
	expectedAuthHeader := `EAN apikey=abc,signature=8333ffbf1c4c4753018554fb91c685cd050cb4c407b473f56fe7405031bb75ffb778f3cd8dc7b71dc4ccae436dc7e6b028bf498912b9bea4197d1f032ac3779b,timestamp=1559215747`

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedAuthHeader, r.Header.Get("Authorization"))
	})

	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := client{url: "http://test_url", Client: httpCli}

	client.getRegions()
}

func TestCreateRequest(t *testing.T) {

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "gzip", r.Header.Get("Accept-Encoding"))
		assert.Equal(t, "BigLife/0.1", r.Header.Get("User-Agent"))
		assert.Equal(t, []string{"details", "property_ids", "property_ids_expanded"}, r.URL.Query()["include"])

	})
	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := client{url: "http://test.com", Client: httpCli}
	client.getRegions()
}

func TestDecode(t *testing.T) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", "regions_stub.txt"))
	if err != nil {
		panic(err)
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()
		_, err := gzipWriter.Write(b)
		err = gzipWriter.Flush()
		if err != nil {
			panic(err)
		}
	})

	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := client{url: "http://test.com", Client: httpCli}

	regions := client.getRegions()
	assert.Equal(t, "Nigeria", regions["136"].Name)
}
