#!/bin/bash

mkdir -p /var/run/nholuongut

for arg in "$@"; do
    case "$arg" in
        --no-app | --probe-only | --service-token* | --probe.token*)
            touch /etc/service/app/down
            ;;
        --no-probe | --app-only)
            touch /etc/service/probe/down
            ;;
    esac
done

# shellcheck disable=SC2034
ARGS=("$@")

declare -p ARGS >/var/run/nholuongut/scope-app.args
# shellcheck disable=SC2034
declare -p ARGS >/var/run/nholuongut/scope-probe.args

exec /home/nholuongut/runsvinit
