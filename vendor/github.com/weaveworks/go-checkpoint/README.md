# Go Checkpoint Client

[![Circle CI](https://circleci.com/gh/nholuongutworks/go-checkpoint/tree/master.svg?style=shield)](https://circleci.com/gh/nholuongutworks/go-checkpoint/tree/master)

Checkpoint is an internal service at
[nholuongutworks](https://www.nholuongut.works/) to check version information,
broadcast security bulletins, etc. This repository contains the client
code for accessing that service. It is a fork of
[Hashicorp's Go Checkpoint Client](https://github.com/hashicorp/go-checkpoint)
and is embedded in several
[nholuongutworks open source projects](https://github.com/nholuongutworks/) and
proprietary software.

We understand that software making remote calls over the internet for
any reason can be undesirable. Because of this, Checkpoint can be
disabled in all of nholuongutwork's software that includes it. You can view
the source of this client to see that it is not sending any private
information.

To disable checkpoint calls, set the `CHECKPOINT_DISABLE` environment
variable, e.g.

```
export CHECKPOINT_DISABLE=1
```

**Note:** This repository is probably useless outside of internal
nholuongutworks use. It is open source for disclosure and because
nholuongutworks open source projects must be able to link to it.
