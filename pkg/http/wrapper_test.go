package http

import (
	"context"
	"io"
	"net/http"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	response := Response{
		StatusCode: 200,
	}

	assert.NoError(t, response.Error())
	assert.True(t, response.IsSuccessful())

	response.StatusCode = 201
	assert.NoError(t, response.Error())
	assert.True(t, response.IsSuccessful())

	response.StatusCode = 400
	assert.Error(t, response.Error())
	assert.False(t, response.IsSuccessful())
}

func TestNewWrapper(t *testing.T) {
	wrapper := NewWrapper()
	assert.NotNil(t, wrapper)
}

func TestNewWrapperRedirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/final" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wrapper := NewWrapper()
	resp, err := wrapper.Client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestWrapperRequestWithError(t *testing.T) {
	wrapper := NewWrapper()

	invalidURL := "http://[::1]:namedport"
	response, err := wrapper.Request(context.Background(), http.MethodGet, invalidURL, Request{})

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestWrapperGet(t *testing.T) {
	mux := http2.NewServeMux()
	mux.HandleFunc("/test", func(w http2.ResponseWriter, r *http2.Request) {
		assert.Equal(t, r.Method, http2.MethodGet)
		_, _ = w.Write([]byte(`test`))
	})

	s := httptest.NewServer(mux)
	defer s.Close()

	wrapper := NewWrapper()

	response, err := wrapper.Get(context.TODO(), s.URL+"/test", Request{})
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http2.StatusOK)
	assert.Equal(t, response.Body, []byte(`test`))
	assert.Nil(t, response.Error())
	assert.True(t, response.IsSuccessful())
}

func TestWrapperPost(t *testing.T) {
	mux := http2.NewServeMux()
	mux.HandleFunc("/test", func(w http2.ResponseWriter, r *http2.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, r.Method, http2.MethodPost)
		assert.Equal(t, r.Header.Get("test"), "test-header")
		assert.Equal(t, body, []byte("test-body"))

		_, _ = w.Write([]byte(`test`))
	})

	s := httptest.NewServer(mux)
	defer s.Close()

	wrapper := NewWrapper()

	response, err := wrapper.Post(context.Background(), s.URL+"/test", Request{
		Headers: map[string]string{
			"test": "test-header",
		},
		Body: []byte("test-body"),
	})
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http2.StatusOK)
	assert.Equal(t, response.Body, []byte(`test`))
}

func TestWrapperPatch(t *testing.T) {
	mux := http2.NewServeMux()
	mux.HandleFunc("/test", func(w http2.ResponseWriter, r *http2.Request) {
		assert.Equal(t, r.Method, http2.MethodPatch)
		_, _ = w.Write([]byte(`test`))
	})

	s := httptest.NewServer(mux)
	defer s.Close()

	wrapper := NewWrapper()

	response, err := wrapper.Patch(context.Background(), s.URL+"/test", Request{})
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http2.StatusOK)
	assert.Equal(t, response.Body, []byte(`test`))
}

func TestWrapperPut(t *testing.T) {
	mux := http2.NewServeMux()
	mux.HandleFunc("/test", func(w http2.ResponseWriter, r *http2.Request) {
		assert.Equal(t, r.Method, http2.MethodPut)
		_, _ = w.Write([]byte(`test`))
	})

	s := httptest.NewServer(mux)
	defer s.Close()

	wrapper := NewWrapper()

	response, err := wrapper.Put(context.Background(), s.URL+"/test", Request{})
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http2.StatusOK)
	assert.Equal(t, response.Body, []byte(`test`))
}

func TestWrapperDelete(t *testing.T) {
	mux := http2.NewServeMux()
	mux.HandleFunc("/test", func(w http2.ResponseWriter, r *http2.Request) {
		assert.Equal(t, r.Method, http2.MethodDelete)
		_, _ = w.Write([]byte(`test`))
	})

	s := httptest.NewServer(mux)
	defer s.Close()

	wrapper := NewWrapper()

	response, err := wrapper.Delete(context.Background(), s.URL+"/test", Request{})
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http2.StatusOK)
	assert.Equal(t, response.Body, []byte(`test`))
}

func TestWrapperRequest(t *testing.T) {
	mux := http2.NewServeMux()
	mux.HandleFunc("/test", func(w http2.ResponseWriter, r *http2.Request) {})

	s := httptest.NewServer(mux)
	defer s.Close()

	wrapper := NewWrapper()

	response, err := wrapper.Request(context.Background(), "invalid", "!invalid-url/test", Request{})
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestWrapperRequestReadError(t *testing.T) {
	mux := http2.NewServeMux()
	mux.HandleFunc("/test", func(w http2.ResponseWriter, r *http2.Request) {
		w.Header().Add("Content-Length", "50")
		_, _ = w.Write([]byte("a"))
	})

	s := httptest.NewServer(mux)
	defer s.Close()

	wrapper := NewWrapper()

	response, err := wrapper.Request(context.Background(), http2.MethodPost, s.URL+"/test", Request{})
	assert.Error(t, err)
	assert.Nil(t, response)
}
