#! /bin/bash

# shellcheck disable=SC1091
. ./config.sh

start_suite "Check scope exits cleanly within 5 seconds"

scope_on "$HOST1" launch
sleep 5
scope_on "$HOST1" stop

sleep 5

# Save stdout for debugging output
exec 3>&1
assert_raises "docker_on $HOST1 logs nholuongutscope 2>&1 | grep 'app exiting' || (docker_on $HOST1 logs nholuongutscope 2>&3 ; false)"
assert_raises "docker_on $HOST1 logs nholuongutscope 2>&1 | grep 'probe exiting' || (docker_on $HOST1 logs nholuongutscope 2>&3 ; false)"
assert_raises "docker_on $HOST1 inspect --format='{{.State.Running}}' nholuongutscope" "false"

scope_end_suite
