package exapi

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-api-samp/util/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExAPISuccess(t *testing.T) {
	t.Run("正常系テスト", func(t *testing.T) {

		server := makeExAPIServer(t)

		apiClient := &MetaWeatherClient{
			Config: &config.ExAPIConfig{
				URL:     server.URL + "/get",
				Timeout: 3 * time.Second,
			},
		}

		result, err := apiClient.GetSample(context.Background())
		if assert.NoError(t, err) {
			assert.NotEmpty(t, result)
		}
	})
}

func TestExAPIError(t *testing.T) {
	tests := []struct {
		name         string
		errURLPrefix string
		urlPath      string
	}{
		{
			name:         "http.NewRequest error",
			errURLPrefix: `\u0001`,
			urlPath:      "/get",
		},
		{
			name:         "http.NewRequest error",
			errURLPrefix: `error`,
			urlPath:      "/get",
		},
		{
			name:    "http status not 200",
			urlPath: "/statuserror",
		},
		{
			name:    "read body error",
			urlPath: "/readbodyerror",
		},
	}

	server := makeExAPIServer(t)

	t.Parallel()
	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			apiClient := &MetaWeatherClient{
				Config: &config.ExAPIConfig{
					URL:     tp.errURLPrefix + server.URL + tp.urlPath,
					Timeout: 3 * time.Second,
				},
			}

			result, err := apiClient.GetSample(context.Background())
			assert.Error(t, err)
			assert.Empty(t, result)
		})
	}
}

func makeExAPIServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		mockHandler(t, writer, request)
	})
	mux.HandleFunc("/statuserror", func(writer http.ResponseWriter, request *http.Request) {
		mockStatusErrorHandler(t, writer, request)
	})
	mux.HandleFunc("/readbodyerror", func(writer http.ResponseWriter, request *http.Request) {
		mockReadBodyErrorHandler(t, writer, request)
	})

	return httptest.NewServer(mux)
}

func mockHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		t.Log("http method is not GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)

	res := fmt.Sprintf(`{"result":%strue}`, "\u0001")
	if _, err := w.Write([]byte(res)); err != nil {
		t.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func mockStatusErrorHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func mockReadBodyErrorHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "1")
}
