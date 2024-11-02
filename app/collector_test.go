package app_test

import (
	"testing"
	"time"

	"context"

	"github.com/nholuongutworks/common/mtime"
	"github.com/nholuongutworks/common/test"
	"github.com/nholuongut/scope/app"
	"github.com/nholuongut/scope/report"
	"github.com/nholuongut/scope/test/reflect"
)

func TestCollector(t *testing.T) {
	ctx := context.Background()
	window := 10 * time.Second
	c := app.NewCollector(window)

	now := time.Now()
	mtime.NowForce(now)
	defer mtime.NowReset()

	r1 := report.MakeReport()
	r1.Endpoint.AddNode(report.MakeNode("foo"))

	r2 := report.MakeReport()
	r2.Endpoint.AddNode(report.MakeNode("foo"))

	have, err := c.Report(ctx, mtime.Now())
	if err != nil {
		t.Error(err)
	}
	if want := report.MakeReport(); !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}

	c.Add(ctx, r1, nil)
	have, err = c.Report(ctx, mtime.Now())
	if err != nil {
		t.Error(err)
	}
	if want := r1; !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}

	timeBefore := mtime.Now()
	mtime.NowForce(now.Add(time.Second))

	c.Add(ctx, r2, nil)
	merged := report.MakeReport()
	merged = merged.Merge(r1)
	merged = merged.Merge(r2)
	have, err = c.Report(ctx, mtime.Now())
	if err != nil {
		t.Error(err)
	}
	if want := merged; !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}

	// Since the timestamp given is before r2 was added,
	// it shouldn't be included in the final report.
	have, err = c.Report(ctx, timeBefore)
	if err != nil {
		t.Error(err)
	}
	if want := r1; !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestCollectorExpire(t *testing.T) {
	now := time.Now()
	mtime.NowForce(now)
	defer mtime.NowReset()

	ctx := context.Background()
	window := 10 * time.Second
	c := app.NewCollector(window)

	// 1st check the collector is empty
	have, err := c.Report(ctx, mtime.Now())
	if err != nil {
		t.Error(err)
	}
	if want := report.MakeReport(); !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}

	// Now check an added report is returned
	r1 := report.MakeReport()
	r1.Endpoint.AddNode(report.MakeNode("foo"))
	c.Add(ctx, r1, nil)
	have, err = c.Report(ctx, mtime.Now())
	if err != nil {
		t.Error(err)
	}
	if want := r1; !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}

	// Finally move time forward to expire the report
	mtime.NowForce(now.Add(window))
	have, err = c.Report(ctx, mtime.Now())
	if err != nil {
		t.Error(err)
	}
	if want := report.MakeReport(); !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestCollectorWait(t *testing.T) {
	ctx := context.Background()
	window := time.Millisecond
	c := app.NewCollector(window)

	waiter := make(chan struct{}, 1)
	c.WaitOn(ctx, waiter)
	defer c.UnWait(ctx, waiter)
	c.(interface {
		Broadcast()
	}).Broadcast()

	select {
	case <-waiter:
	default:
		t.Fatal("Didn't unblock")
	}
}
