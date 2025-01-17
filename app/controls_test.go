package app_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/ugorji/go/codec"

	"github.com/nholuongut/scope/app"
	"github.com/nholuongut/scope/common/xfer"
	"github.com/nholuongut/scope/probe/appclient"
)

func TestControl(t *testing.T) {
	router := mux.NewRouter()
	app.RegisterControlRoutes(router, app.NewLocalControlRouter())
	server := httptest.NewServer(router)
	defer server.Close()

	ip, port, err := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
	if err != nil {
		t.Fatal(err)
	}

	probeConfig := appclient.ProbeConfig{
		ProbeID: "foo",
	}
	controlHandler := xfer.ControlHandlerFunc(func(req xfer.Request) xfer.Response {
		if req.NodeID != "nodeid" {
			t.Fatalf("'%s' != 'nodeid'", req.NodeID)
		}

		if req.Control != "control" {
			t.Fatalf("'%s' != 'control'", req.Control)
		}

		return xfer.Response{
			Value: "foo",
		}
	})
	url := url.URL{Scheme: "http", Host: ip + ":" + port}
	client, err := appclient.NewAppClient(probeConfig, ip+":"+port, url, controlHandler)
	if err != nil {
		t.Fatal(err)
	}
	client.ControlConnection()
	defer client.Stop()

	time.Sleep(100 * time.Millisecond)

	httpClient := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := httpClient.Post(
		server.URL+"/api/control/foo/nodeid/control",
		"application/json",
		strings.NewReader("{}"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var response xfer.Response
	decoder := codec.NewDecoder(resp.Body, &codec.JsonHandle{})
	if err := decoder.Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Value != "foo" {
		t.Fatalf("'%s' != 'foo'", response.Value)
	}
}
