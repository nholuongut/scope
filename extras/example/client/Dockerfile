FROM tatsushid/tinycore-python:2.7
WORKDIR /home/nholuongut
ADD requirements.txt /home/nholuongut/
RUN pip install -r /home/nholuongut/requirements.txt
ADD client.py /home/nholuongut/
ENTRYPOINT ["python", "/home/nholuongut/client.py"]

ARG revision
LABEL maintainer="nholuongutworks <luongutnho@hotmail.com>" \
      org.opencontainers.image.title="example-client" \
      org.opencontainers.image.source="https://github.com/nholuongut/scope/tree/master/extras/example/client" \
      org.opencontainers.image.revision="${revision}" \
      org.opencontainers.image.vendor="nholuongutworks"
