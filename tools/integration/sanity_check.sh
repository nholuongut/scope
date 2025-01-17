#! /bin/bash
# shellcheck disable=SC1090,SC1091
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/config.sh"

set -e

whitely echo Ping each host from the other
for host in $HOSTS; do
    for other in $HOSTS; do
        [ "$host" = "$other" ] || run_on "$host" "$PING" "$other"
    done
done

whitely echo Check we can reach docker

for host in $HOSTS; do
    echo
    echo "Host Version Info: $host"
    echo "====================================="
    echo "# docker version"
    docker_on "$host" version
    echo "# docker info"
    docker_on "$host" info
    echo "# nholuongut version"
    nholuongut_on "$host" version
done
