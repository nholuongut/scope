FROM progrium/busybox
WORKDIR /home/nholuongut
ADD searchapp /home/nholuongut/
EXPOSE 8080
ENTRYPOINT ["/home/nholuongut/searchapp"]

ARG revision
LABEL maintainer="nholuongutworks <luongutnho@hotmail.com>" \
      org.opencontainers.image.title="example-searchapp" \
      org.opencontainers.image.source="https://github.com/nholuongut/scope/tree/master/extras/example/searchapp" \
      org.opencontainers.image.revision="${revision}" \
      org.opencontainers.image.vendor="nholuongutworks"
