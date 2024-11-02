#! /bin/bash

# shellcheck disable=SC1091
. ./config.sh

start_suite "Launch scope and check it boots, with a spurious host arg"

scope_on "$HOST1" launch noatrealhost.foo

wait_for_containers "$HOST1" 60 nholuongutscope

has_container "$HOST1" nholuongutscope

scope_end_suite
