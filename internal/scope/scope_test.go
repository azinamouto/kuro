package scope

import "testing"

func TestScope(t *testing.T) {
	p := Policy{Version: 1, Origin: "http://127.0.0.1:8080", MaxRequests: 10}
	for _, tc := range []struct {
		url string
		ok  bool
	}{{"http://127.0.0.1:8080/api", true}, {"http://127.0.0.1:8081", false}, {"http://localhost:8080", false}, {"http://192.168.1.1:8080", false}, {"http://127.0.0.1:8080@evil.test/", false}, {"https://127.0.0.1:8080", false}, {"http://127.0.0.1:8080/a%2fb", false}} {
		err := p.Check(tc.url)
		if (err == nil) != tc.ok {
			t.Errorf("%s: got %v", tc.url, err)
		}
	}
}

