#! /bin/bash

# shellcheck disable=SC1091
. ./config.sh

start_suite "Launch scope and check it boots"

nholuongut_on "$HOST1" launch
scope_on "$HOST1" launch

wait_for_containers "$HOST1" 60 nholuongut nholuongutscope

has_container "$HOST1" nholuongut
has_container "$HOST1" nholuongutscope

# Fail if the top-level UI is suspiciously small
ui_len="$(curl -s "http://$HOST1:4040/" | wc -c)"
assert_raises "(( $ui_len > 500 ))" 0

scope_end_suite
