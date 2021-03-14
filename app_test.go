package main

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

const (
	expTime        = 60
	longURL        = "https://www.qq.com"
	shortLink      = "IJiwa"
	shortLinkInfor = `{"url": "https://www.qq.com", "created_at":"2018-03-21 15:33:13"`
)

type storageMock struct {
	mock.Mock
}

var app App
var mockR *storageMock

func (s *storageMock) Shorten(url string, exp int64) (string, error) {
	args := s.Called(url, exp)
	return args.String(0), args.Error(1)
}

func (s *storageMock) Unshorten(eid string) (string, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func (s *storageMock) ShortlinkInfo(eid string) (interface{}, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func init() {
	app = App{}
	mockR = new(storageMock)
	app.Initialize(&Env{S: mockR})
}

func TestCreateShortlink(t *testing.T) {
	var jsonStr = []byte(`{
		"url": "https://www.qq.com"
		"expiration_in_minutes": 60
		}`)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal("during create a request")
	}
	
	req.Header.Set("Content-Type", "application/json")
	mockR.On("Shorten", longURL, int64(expTime)).Return(shortLink, nil).Once()
	rw := httptest.NewRecorder()

	app.Router.ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("Excepted receive %d. Got %d", http.StatusCreated, rw.Code)
	}

	resp := struct {
		Shortlink string `json:"shortlink"`
	}{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("during decode the response")
	}

	if resp.Shortlink != shortLink {
		t.Fatalf("during receive %s. Got %s", shortLink, resp.Shortlink)
	}

}
