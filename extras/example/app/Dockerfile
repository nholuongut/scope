FROM tatsushid/tinycore-python:2.7
WORKDIR /home/nholuongut
ADD requirements.txt /home/nholuongut/
RUN pip install -r /home/nholuongut/requirements.txt
ADD app.py /home/nholuongut/
EXPOSE 5000
ENTRYPOINT ["python", "/home/nholuongut/app.py"]

ARG revision
LABEL maintainer="nholuongutworks <luongutnho@hotmail.com>" \
      org.opencontainers.image.title="example-app" \
      org.opencontainers.image.source="https://github.com/nholuongut/scope/tree/master/extras/example/app" \
      org.opencontainers.image.revision="${revision}" \
      org.opencontainers.image.vendor="nholuongutworks"
