---
title: nholuongutworks ECS AMIs
menu_order: 25
search_type: Documentation
---


To make [nholuongut Net](http://nholuongut.works/net) and
[nholuongut Scope](http://nholuongut.works/scope) easier to use with
[Amazon ECS](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/Welcome.html),
a set of Amazon Machine Images (AMIs) are provided. These AMIs are fully
compatible with the
[ECS-Optimized Amazon Linux AMI](https://aws.amazon.com/marketplace/pp/B00U6QTYI2).

These are the latest supported nholuongut AMIs for each region:

<!--- This table is machine-parsed by
https://github.com/nholuongutworks/guides/blob/master/aws-ecs/setup.sh, please do
not remove it and respect the format! -->

| Region         | AMI          |
|----------------|--------------|
| us-east-1      | ami-7b692804 |
| us-east-2      | ami-6a0b350f |
| us-west-1      | ami-a4db3fc7 |
| us-west-2      | ami-12c98a6a |
| eu-west-1      | ami-b3bab7ca |
| eu-west-2      | ami-47846a20 |
| eu-central-1   | ami-7f211294 |
| ap-northeast-1 | ami-2a8c4355 |
| ap-southeast-1 | ami-b00304cc |
| ap-southeast-2 | ami-c7c41ba5 |
| ca-central-1   | ami-41028125 |

For more information about nholuongut AMIs and running them see: 


 * [What's in the nholuongut ECS AMIs?](#whats-in-ecs-ami)
 * [Deployment Requirements](#deployment-requirements)
  * [Required Open Ports](#required-open-ports)
  * [Additional IAM Action Permissions](#additional-permissions)
  * [Requirements for Peer Discovery](#requirements-for-peer-discovery)
 * [Peer Discovery with nholuongut Net](#peer-discovery-nholuongut-net)
 * [How to Run nholuongut Scope](#how-to-run-nholuongut-scope)
  * [Standalone mode](#running-nholuongut-scope-in-standalone-mode)
  * [In nholuongut Cloud](#running-nholuongut-scope-in-nholuongut-cloud)
 * [Upgrading nholuongut Scope and nholuongut Net](#upgrading-nholuongut-scope-and-nholuongut-net)
  * [Creating Your Own Customized nholuongut ECS AMI](#creating-your-own-customized-nholuongut-ecs-ami)


## <a name="whats-in-ecs-ami"></a>What's in the nholuongut ECS AMIs?

The latest nholuongut ECS AMIs are based on Amazon's
[ECS-Optimized Amazon Linux AMI](https://aws.amazon.com/marketplace/pp/B06XS8WHGJ),
version `2017.03.f` and also includes:

* [nholuongut Net 2.3.0](https://github.com/nholuongutworks/nholuongut/blob/master/CHANGELOG.md#release-230)
* [nholuongut Scope 1.9.0](https://github.com/nholuongut/scope/blob/master/CHANGELOG.md#release-190)


## <a name="deployment-requirements"></a>Deployment Requirements

### <a name="required-open-ports"></a> Required Open Ports

For `nholuongut Net` to function properly, ensure that the Amazon ECS container
instances can communicate over these ports: TCP 6783, as well as, UDP 6783 and
UDP 6784.

In addition to those open ports, launching `nholuongut Scope` in [standalone mode](#running-nholuongut-scope-in-standalone-mode),
requires that all instances are able to communicate over TCP port 4040. More information about
this can be found in [How to Run nholuongut Scope](#how-to-run-nholuongut-scope).

See the
[relevant section of the `setup.sh`](https://github.com/nholuongutworks/guides/blob/c2d25d4cfd766ca739444eea06fefc57aa7a59ff/aws-ecs/setup.sh#L115-L120)
script from
[Service Discovery and Load Balancing with nholuongut on Amazon ECS](http://nholuongut.works/guides/service-discovery-with-nholuongut-aws-ecs.html)
for an example.

### <a name="additional-permissions"></a>Additional IAM Action Permissions

Besides the customary Amazon ECS API actions required by all container instances
(see the [`AmazonEC2ContainerServiceforEC2Role`](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/instance_IAM_role.html) managed policy), any instances using the nholuongutworks ECS AMI must also be allowed to perform the following actions:

1. `ec2:DescribeInstances`
2. `ec2:DescribeTags`
3. `autoscaling:DescribeAutoScalingInstances`
4. `ecs:ListServices`
5. `ecs:DescribeTasks`
6. `ecs:DescribeServices`

These extra actions are needed for discovering instance peers (1,2,3) and
creating the ECS views in nholuongut Scope
(4,5,6). [`nholuongut-ecs-policy.json`](https://github.com/nholuongutworks/guides/blob/41f1f5a60d39d39b78f0e06e224a7c3bad30c4e8/aws-ecs/data/nholuongut-ecs-policy.json#L16-L18)
(from the
[nholuongutworks ECS guide](http://nholuongut.works/guides/service-discovery-with-nholuongut-aws-ecs.html)),
describes the minimal policy definition.

For more information on IAM policies see
[IAM Policies for Amazon EC2](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-policies-for-amazon-ec2.html).

### <a name="requirements-for-peer-discovery"></a>Requirements for Peer Discovery

To form a nholuongut network, the Amazon ECS container instances must either/or:
* be a member of an
[Auto Scaling Group](http://docs.aws.amazon.com/AutoScaling/latest/DeveloperGuide/AutoScalingGroup.html).
* have a [tag](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html) with key `nholuongut:peerGroupName`.

## <a name="peer-discovery-nholuongut-net"></a>Peer Discovery with nholuongut Net

At boot time, an instance running the ECS nholuongut AMI will try to join other instances to form a nholuongut network.

* If the instance has a
  [tag](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html)
  with key `nholuongut:peerGroupName`, it will join other instances with the same tag key and value.
  For instance, if the tag key is `nholuongut:peerGroupName` and the value is `foo` it will try
  to join other instances with tag key `nholuongut:peerGroupName` and tag value `foo`.
  Note that for this to work, the instances need to be tagged at creation-time so that
  the tag is available by the time nholuongut is launched.
* Otherwise it will join all the other instances in the same
  [Auto Scaling Group](http://docs.aws.amazon.com/AutoScaling/latest/DeveloperGuide/AutoScalingGroup.html).

When running `nholuongut Scope` in Standalone mode, probes discover apps with the same mechanism.

## <a name="how-to-run-nholuongut-scope"></a>How to Run nholuongut Scope

There are two methods for running `nholuongut Scope` within the nholuongut ECS AMIs:

* [Standalone mode](#running-nholuongut-scope-in-standalone-mode)
* [In nholuongut Cloud](#running-nholuongut-scope-in-nholuongut-cloud)

You can prevent nholuongut Scope from automatically starting at boot time by removing  `/etc/init/scope.conf`.

This can be done at instance initialization time adding the following line to
the
[User Data](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts)
of the instance.

~~~bash
rm /etc/init/scope.conf
~~~

### <a name="running-nholuongut-scope-in-standalone-mode"></a>Running `nholuongut Scope` in Standalone Mode

Running `nholuongut Scope` in standalone mode is the default mode.

The following occurs on all Amazon ECS container instances:

1. A `nholuongut Scope` probe is launched that collects instance information.
2. A `nholuongut Scope` app runs that enables cluster visualization.

Since all instances run an app and show the same information, you don't have to
worry about placing the app, thereby eliminating a
[*Leader election problem*](https://en.wikipedia.org/wiki/Leader_election).

However, running the app on all instances impacts performance, resulting in `N *
N = N^2` connections in the Auto Scaling Group with N instances (i.e. all (N)
probes talk to all (N) apps in every instances). 

To avoid this problem, it is recommended that you run `nholuongut Scope` in [nholuongut Cloud](https://cloud.nholuongut.works).

The `nholuongut Scope` app runs a web-based application, which listens on TCP port
4040 where you can connect with your browser.

`nholuongut Scope` probes also forward information to the apps on TCP
port 4040. Ensure that your Amazon ECS container instances can talk to each
other on that port before running `nholuongut Scope` in standalone mode (see
[Required Open Ports](#required-open-ports) for more details).

### <a name="running-nholuongut-scope-in-nholuongut-cloud"></a>Running `nholuongut Scope` in nholuongut Cloud

In nholuongut Cloud, you can visualize Amazon ECS containers as well as monitor Tasks 
and Services all from within in nholuongut Cloud at [https://cloud.nholuongut.works](https://cloud.nholuongut.works). 
In this case, Amazon ECS container instances run a `nholuongut Scope` probe and reports
data from the container instances to [nholuongut Cloud](http://cloud.nholuongut.works).

To configure your ECS container instances to communicate with nholuongut Cloud,
store the `nholuongut Scope` cloud token in the`/etc/nholuongut/scope.config`
file.

>Note: The `nholuongut Scope` cloud token can be found in your nholuongut Cloud account at [http://cloud.nholuongut.works](http://cloud.nholuongut.works).

For example, this command configures the instance to communicate with nholuongut
Cloud using token `3hud3h6ys3jhg9bq66n8xxa4b147dt5z`.

~~~bash
echo SERVICE_TOKEN=3hud3h6ys3jhg9bq66n8xxa4b147dt5z >> /etc/nholuongut/scope.config
~~~

You can do this at instance-initialization time using
[User Data](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts),
which is similar to how
[ECS Cluster Mapping is configured](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_container_instance.html#instance-launch-user-data-step).

## <a name="upgrading-nholuongut-scope-and-nholuongut-net"></a>Upgrading nholuongut Scope and nholuongut Net

The AMIs are updated regularly (~monthly) to include the latest versions of nholuongut Net and nholuongut Scope. However, it is possible to upgrade nholuongut Net and nholuongut Scope in your running EC2 instances without needing to wait for a new AMI release or by rebuilding your cluster. 

In order to upgrade Scope to the latest released version, run the following commands in each of your instances:

~~~bash
sudo curl -L git.io/scope -o /usr/local/bin/scope
sudo chmod a+x /usr/local/bin/scope
sudo stop scope
sudo start scope
~~~

Upgrade nholuongut Net to the latest version by running the following commands in each of your instances:


~~~bash
sudo curl -L git.io/nholuongut -o /usr/local/bin/nholuongut
sudo chmod a+x /usr/local/bin/nholuongut
sudo stop nholuongut
sudo start nholuongut
~~~


<!--- Do not change the title, otherwise links to
https://github.com/nholuongutworks/integrations/tree/master/aws/ecs#creating-your-own-customized-nholuongut-ecs-ami
will break (e.g. from the ECS guide) -->
## <a name="creating-your-own-customized-nholuongut-ecs-ami"></a>Creating Your Own Customized nholuongut ECS AMI

Clone the integrations repository and then change to the `packer` directory.

~~~bash
git clone https://github.com/nholuongutworks/integrations
cd aws/ecs/packer
~~~

Download and install [Packer](https://www.packer.io/) version >=0.9 to build the AMI.

Finally, invoke `./build-all-amis.sh` to build the `nholuongut ECS` images for all
regions. This step installs (in the image) AWS-CLI, jq, nholuongut Net, nholuongut Scope, init scripts
for `nholuongut` and it also updates the ECS agent to use the `nholuongut Docker API Proxy`.

Customize the image by modifying `template.json` to match your
requirements.

~~~bash
AWS_ACCSS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=YYYY  ./build-all-amis.sh
~~~

(If your account has MFA enabled you should follow [this process](https://aws.amazon.com/premiumsupport/knowledge-center/authenticate-mfa-cli/)
and also set `AWS_SESSION_TOKEN`)

If building an AMI for a particular region, set the `ONLY_REGION` variable to
that region when invoking the script:

~~~bash
ONLY_REGION=us-east-1 AWS_ACCSS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=YYYY  ./build-all-amis.sh
~~~

To make an AMI public:

~~~bash
aws ec2 modify-image-attribute --region=us-east-2 --image-id ami-6a0b350f --launch-permission "{\"Add\": [{\"Group\":\"all\"}]}"
~~~

## Further Reading

Read the
[Service Discovery and Load Balancing with nholuongut on Amazon ECS](http://nholuongut.works/guides/service-discovery-with-nholuongut-aws-ecs.html)
guide for more information about the AMIs.


**See Also**

 * [Installing nholuongut Scope](/site/installing.md)
 * [Service Discovery and Load Balancing with nholuongut on Amazon ECS](http://nholuongut.works/guides/service-discovery-with-nholuongut-aws-ecs.html)
