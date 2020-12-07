package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/creativeprojects/clog"
	"github.com/stretchr/testify/assert"
)

func TestSimpleStream(t *testing.T) {
	// test HTTP/1.1
	t.Run("HTTP1", func(t *testing.T) {
		server := httptest.NewServer(getServeMux())
		testSimpleStream(t, server)
	})
	// test HTTP/2.0
	t.Run("HTTP2", func(t *testing.T) {
		server := httptest.NewUnstartedServer(getServeMux())
		server.EnableHTTP2 = true
		server.StartTLS()
		testSimpleStream(t, server)
	})
}

func testSimpleStream(t *testing.T, server *httptest.Server) {
	testData := []struct {
		name    string
		waitPut time.Duration
		waitGet time.Duration
	}{
		{"PutFirst", 10 * time.Millisecond, 20 * time.Millisecond},
		{"GetFirst", 20 * time.Millisecond, 10 * time.Millisecond},
		{"NoWait", 0, 0},
	}

	client := server.Client()
	t.Logf("test server listening on %s", server.URL)

	for _, testRun := range testData {
		t.Run(testRun.name, func(t *testing.T) {
			clog.SetTestLog(t)
			defer clog.CloseTestLog()

			messagePut := fmt.Sprintf("Hello TestSimpleStream %s!", testRun.name)
			messageGet := ""
			url := server.URL + "/TestSimpleStream/" + testRun.name
			wg := &sync.WaitGroup{}

			// send message
			wg.Add(1)
			go func() {
				time.Sleep(testRun.waitPut)
				req, err := http.NewRequest("PUT", url, bytes.NewReader([]byte(messagePut)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "text/plain")
				resp, err := client.Do(req)
				assert.NoError(t, err)

				// drain the answer (should be empty)
				buffer := &bytes.Buffer{}
				_, err = io.Copy(buffer, resp.Body)
				assert.NoError(t, err)

				err = resp.Body.Close()
				assert.NoError(t, err)
				wg.Done()
			}()

			// read message back
			wg.Add(1)
			go func() {
				time.Sleep(testRun.waitGet)
				resp, err := client.Get(url)
				assert.NoError(t, err)

				t.Log(resp.Proto)

				buffer := &strings.Builder{}
				_, err = io.Copy(buffer, resp.Body)
				assert.NoError(t, err)

				err = resp.Body.Close()
				assert.NoError(t, err)

				messageGet = buffer.String()
				wg.Done()
			}()

			wg.Wait()
			assert.Equal(t, messagePut, messageGet)
		})
	}
	server.Close()
}

func TestNoPOST(t *testing.T) {
	clog.SetTestLog(t)
	defer clog.CloseTestLog()

	server := httptest.NewServer(getServeMux())
	client := server.Client()
	req, err := client.Post(server.URL+"/TestNoPOST", "text/plain", bytes.NewReader([]byte("message")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, req.StatusCode)
}

func TestTwoGetBeforePut(t *testing.T) {
	server := httptest.NewServer(getServeMux())
	client := server.Client()

	clog.SetTestLog(t)
	defer clog.CloseTestLog()

	messagePut := fmt.Sprintf("Hello TestTwoGetBeforePut!")
	messageGet := ""
	url := server.URL + "/TestTwoGetBeforePut"
	wg := &sync.WaitGroup{}

	// send message
	wg.Add(1)
	go func() {
		time.Sleep(10 * time.Millisecond)
		req, err := http.NewRequest("PUT", url, bytes.NewReader([]byte(messagePut)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")
		resp, err := client.Do(req)
		assert.NoError(t, err)

		// drain the answer (should be empty)
		buffer := &bytes.Buffer{}
		_, err = io.Copy(buffer, resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)
		wg.Done()
	}()

	// 1st read
	wg.Add(1)
	go func() {
		time.Sleep(1 * time.Millisecond)
		resp, err := client.Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		buffer := &strings.Builder{}
		_, err = io.Copy(buffer, resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		messageGet = buffer.String()
		wg.Done()
	}()

	// 2nd read
	wg.Add(1)
	go func() {
		time.Sleep(3 * time.Millisecond)
		resp, err := client.Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		err = resp.Body.Close()
		assert.NoError(t, err)

		wg.Done()
	}()

	wg.Wait()
	assert.Equal(t, messagePut, messageGet)

	server.Close()
}
