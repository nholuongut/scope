#! /bin/bash

# shellcheck disable=SC1091
. ./config.sh

start_suite "Launch scope (without nholuongut installed) and check it boots"

scope_on "$HOST1" launch

wait_for_containers "$HOST1" 60 nholuongutscope

has_container "$HOST1" nholuongut 0
has_container "$HOST1" nholuongutscope

scope_end_suite
