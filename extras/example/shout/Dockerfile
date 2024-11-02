FROM alpine:latest
WORKDIR /home/nholuongut
ADD shout /home/nholuongut/
EXPOSE 8090
ENTRYPOINT ["/home/nholuongut/shout"]

ARG revision
LABEL maintainer="nholuongutworks <luongutnho@hotmail.com>" \
      org.opencontainers.image.title="example-shout" \
      org.opencontainers.image.source="https://github.com/nholuongut/scope/tree/master/extras/example/shout" \
      org.opencontainers.image.revision="${revision}" \
      org.opencontainers.image.vendor="nholuongutworks"
