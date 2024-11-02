FROM gliderlabs/alpine
WORKDIR /
COPY proxy /
EXPOSE 8000
EXPOSE 8080
ENTRYPOINT ["/proxy"]

ARG revision
LABEL maintainer="nholuongutworks <luongutnho@hotmail.com>" \
      org.opencontainers.image.title="socks" \
      org.opencontainers.image.source="https://github.com/nholuongutworks/build-tools/tree/master/socks" \
      org.opencontainers.image.revision="${revision}" \
      org.opencontainers.image.vendor="nholuongutworks"
