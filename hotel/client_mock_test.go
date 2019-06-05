package hotel

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"net"
	"net/http"
	"net/http/httptest"
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) getRegions() Regions {
	fmt.Println("mockClient getRegions called")
	args := m.Called()
	fmt.Println("args extracted are :", args[0])
	return args[0].(Regions)
}

func MockHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}
