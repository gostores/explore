package cachecontrol

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gostores/explore/cacheobject"
	"github.com/gostores/require"
)

func roundTrip(t *testing.T, fnc func(w http.ResponseWriter, r *http.Request)) (*http.Request, *http.Response) {
	ts := httptest.NewServer(http.HandlerFunc(fnc))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	return req, res
}

func TestCachableResponsePublic(t *testing.T) {
	req, res := roundTrip(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public")
		w.Header().Set("Last-Modified",
			time.Now().UTC().Add(time.Duration(time.Hour*-5)).Format(http.TimeFormat))
		fmt.Fprintln(w, `{}`)
	})

	opts := Options{}
	reasons, expires, err := CachableResponse(req, res, opts)
	require.NoError(t, err)
	require.Len(t, reasons, 0)
	require.WithinDuration(t,
		time.Now().UTC().Add(time.Duration(float64(time.Hour)*0.5)),
		expires,
		10*time.Second)
}

func TestCachableResponsePrivate(t *testing.T) {
	req, res := roundTrip(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "private")
		fmt.Fprintln(w, `{}`)
	})

	opts := Options{}
	reasons, expires, err := CachableResponse(req, res, opts)
	require.NoError(t, err)
	require.Len(t, reasons, 1)
	require.Equal(t, reasons[0], cacheobject.ReasonResponsePrivate)
	require.Equal(t, time.Time{}, expires)

	opts.PrivateCache = true
	reasons, expires, err = CachableResponse(req, res, opts)
	require.NoError(t, err)
	require.Len(t, reasons, 0)
	require.Equal(t, time.Time{}, expires)
}

func TestResponseWriter(t *testing.T) {
	var resp http.ResponseWriter
	var req *http.Request
	_, _ = roundTrip(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "private")
		fmt.Fprintln(w, `{}`)
		resp = w
		req = r
	})

	opts := Options{}
	reasons, expires, err := CachableResponseWriter(req, 200, resp, opts)
	require.NoError(t, err)
	require.Len(t, reasons, 1)
	require.Equal(t, reasons[0], cacheobject.ReasonResponsePrivate)
	require.Equal(t, time.Time{}, expires)

	opts.PrivateCache = true
	reasons, expires, err = CachableResponseWriter(req, 200, resp, opts)
	require.NoError(t, err)
	require.Len(t, reasons, 0)
	require.Equal(t, time.Time{}, expires)
}
