#!/bin/sh

set -eu

SCOPE_SRC=$GOPATH/src/github.com/nholuongut/scope

# Mount the scope repo:
#  -v $(pwd):/go/src/github.com/nholuongut/scope

# If we run make directly, any files created on the bind mount
# will have awkward ownership.  So we switch to a user with the
# same user and group IDs as source directory.  We have to set a
# few things up so that sudo works without complaining later on.
uid=$(stat --format="%u" "$SCOPE_SRC")
gid=$(stat --format="%g" "$SCOPE_SRC")
echo "nholuongut:x:$uid:$gid::$SCOPE_SRC:/bin/sh" >>/etc/passwd
echo "nholuongut:*:::::::" >>/etc/shadow
echo "nholuongut	ALL=(ALL)	NOPASSWD: ALL" >>/etc/sudoers

su nholuongut -c "PATH=$PATH make -C $SCOPE_SRC BUILD_IN_CONTAINER=false $*"
