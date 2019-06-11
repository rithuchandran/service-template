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

func TestGetRegionsShouldReturnValidRegion(t *testing.T) {
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
	client := NewClient("http://test.com")
	client.Client = httpCli
	client.Timeout = time.Duration(1) * time.Second

	regions, err := client.getRegions()

	assert.Nil(t, err)
	assert.Equal(t, 250, len(regions))
	assert.Equal(t, "Albania", regions["2"].Name)
}

func TestGetRegionsShouldReturnDoError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("test"))
	})
	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := client{url: "http://test.com", Client: httpCli}
	client.Timeout = time.Duration(1) * time.Second

	regions, err := client.getRegions()

	assert.Equal(t, Regions{}, regions)
	assert.EqualError(t, err, "Do error : <nil> ", "expected do error")
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
	client := client{url: "http://test.com", Client: httpCli}
	client.Timeout = time.Duration(1) * time.Second

	client.getRegions()
}

func TestGetNextLink(t *testing.T) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", "regions_stub.txt"))
	if err != nil {
		panic(err)
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Link", "://next_link")
		_, _ = w.Write(b)
	})
	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := NewClient("http://test.com")
	client.Client = httpCli
	client.Timeout = time.Duration(1) * time.Second

	regions, err := client.getRegions()

	assert.EqualError(t, err, "parse ://next_link: missing protocol scheme")
	assert.Equal(t, Regions{}, regions)

}

func TestCreateRequestShouldReturnAValidRequestWithHeaders(t *testing.T) {

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "gzip", r.Header.Get("Accept-Encoding"))
		assert.Equal(t, "BigLife/0.1", r.Header.Get("User-Agent"))
		assert.Equal(t, []string{"details", "property_ids", "property_ids_expanded"}, r.URL.Query()["include"])

	})
	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := client{url: "http://test.com", Client: httpCli}
	client.Timeout = time.Duration(1) * time.Second

	client.getRegions()
}

func TestCreateRequestShouldReturnError(t *testing.T) {
	client := NewClient("://test.com")

	regions, err := client.getRegions()
	assert.Equal(t, Regions{}, regions)
	assert.EqualError(t, err, "parse ://test.com/regions: missing protocol scheme")
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
		assert.NoError(t, err)
	})

	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := client{url: "http://test.com", Client: httpCli}
	client.Timeout = time.Duration(1) * time.Second


	regions, err := client.getRegions()
	assert.Nil(t, err)
	assert.Equal(t, "Nigeria", regions["136"].Name)
}

func TestDecodeShouldReturnGzipError(t *testing.T) {

	b := []byte(`{test data}`)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		_, err := w.Write(b)
		assert.NoError(t, err)
	})

	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := client{url: "http://test", Client: httpCli}
	client.Timeout = time.Duration(1) * time.Second

	regions, err := client.getRegions()
	assert.EqualError(t, err, "gzip: invalid header", "expected gzip error")
	assert.Equal(t, Regions{}, regions)

}

func TestDecodeShouldReturnJsonError(t *testing.T) {

	b := []byte(`{"test data" : "something"}`)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()
		_, err := gzipWriter.Write(b)
		err = gzipWriter.Flush()
		assert.NoError(t, err)
	})

	httpCli, stop := MockHTTPClient(h)
	defer stop()
	client := client{url: "http://test", Client: httpCli}
	client.Timeout = time.Duration(1) * time.Second

	regions, err := client.getRegions()
	assert.EqualError(t, err, "json: cannot unmarshal string into Go value of type hotel.Region", "expected gzip error")
	assert.Equal(t, Regions{}, regions)

}
