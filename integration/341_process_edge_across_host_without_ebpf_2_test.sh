#! /bin/bash

# shellcheck disable=SC1091
. ./config.sh

start_suite "Test long connections (procspy) between processes on different hosts"

nholuongut_on "$HOST1" launch "$HOST1" "$HOST2"
nholuongut_on "$HOST2" launch "$HOST1" "$HOST2"

scope_on "$HOST1" launch --probe.ebpf.connections=false
scope_on "$HOST2" launch --probe.ebpf.connections=false

server_on "$HOST1"
nholuongut_proxy_on "$HOST2" run -dti --name client alpine /bin/sh -c "while true; do \
	nc nginx.nholuongut.local 80 || true; \
	sleep 1; \
done"

sleep 30 # need to allow the scopes to poll dns, resolve the other app ids, and send them reports

check() {
    has processes "$1" "nginx: worker process"
    has processes "$1" nc
    has_connection processes "$1" nc "nginx: worker process"
}

check "$HOST1"
check "$HOST2"

scope_end_suite
