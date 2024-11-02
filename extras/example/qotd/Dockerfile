FROM ubuntu
WORKDIR /home/nholuongut
ADD ./qotd /home/nholuongut/
EXPOSE 4446
ENTRYPOINT ["/home/nholuongut/qotd"]

ARG revision
LABEL maintainer="nholuongutworks <luongutnho@hotmail.com>" \
      org.opencontainers.image.title="example-qotd" \
      org.opencontainers.image.source="https://github.com/nholuongut/scope/tree/master/extras/example/qotd" \
      org.opencontainers.image.revision="${revision}" \
      org.opencontainers.image.vendor="nholuongutworks"
