// Publish a fixed report.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"github.com/ugorji/go/codec"

	"github.com/nholuongut/scope/common/xfer"
	"github.com/nholuongut/scope/probe/appclient"
	"github.com/nholuongut/scope/report"
	"github.com/nholuongut/scope/test/fixture"
)

func main() {
	var (
		publish         = flag.String("publish", fmt.Sprintf("127.0.0.1:%d", xfer.AppPort), "publish target")
		publishInterval = flag.Duration("publish.interval", 1*time.Second, "publish (output) interval")
		publishToken    = flag.String("publish.token", "fixprobe", "publish token, for if we are talking to the service")
		publishID       = flag.String("publish.id", "fixprobe", "publisher ID used to identify publishers")
		useFixture      = flag.Bool("fixture", false, "Use the embedded fixture report.")
	)
	flag.Parse()

	if len(flag.Args()) != 1 && !*useFixture {
		log.Fatal("usage: fixprobe [--args] report.json")
	}

	var fixedReport report.Report
	if *useFixture {
		fixedReport = fixture.Report
	} else {
		b, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		decoder := codec.NewDecoderBytes(b, &codec.JsonHandle{})
		if err := decoder.Decode(&fixedReport); err != nil {
			log.Fatal(err)
		}
	}

	url, err := url.Parse(fmt.Sprintf("http://%s", *publish))
	if err != nil {
		log.Fatal(err)
	}

	client, err := appclient.NewAppClient(appclient.ProbeConfig{
		Token:    *publishToken,
		ProbeID:  *publishID,
		Insecure: false,
	}, *publish, *url, nil)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := fixedReport.WriteBinary()
	if err != nil {
		log.Fatal(err)
	}
	for range time.Tick(*publishInterval) {
		client.Publish(bytes.NewReader(buf.Bytes()), fixedReport.Shortcut)
	}
}
