# example

To run a on a Mac, run:

```
$ docker-machine create -d virtualbox --virtualbox-memory=4096 scope-tastic
$ eval $(docker-machine env scope-tastic)
$ sudo curl -L git.io/nholuongut -o /usr/local/bin/nholuongut
$ sudo chmod +x /usr/local/bin/nholuongut
$ curl -o run.sh https://raw.githubusercontent.com/nholuongut/scope/master/experimental/example/run.sh
$ ./run.sh
$ sudo curl -L git.io/scope -o /usr/local/bin/scope
$ sudo chmod a+x /usr/local/bin/scope
$ scope launch
```

# "architecture"

```
client -> frontend --> app --> searchapp -> elasticsearch
           (nginx)         |
                           --> qotd -> internet
                           |
                           --> redis
```

# To push new images

```
for img in $(docker images | grep tomwilkie | cut -d' ' -f1); do docker push $img:latest; done
```