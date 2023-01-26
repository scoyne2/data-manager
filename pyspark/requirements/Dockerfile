FROM --platform=linux/amd64 amazonlinux:2 AS base

RUN yum install -y python3

ENV VIRTUAL_ENV=/opt/venv
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

COPY ./requirements.txt /
RUN python3 -m pip install --upgrade pip && \
    python3 -m pip install -r requirements.txt --no-cache-dir

RUN python3 -m pip install venv-pack

RUN mkdir /output && venv-pack -o /output/pyspark_requirements.tar.gz

FROM scratch AS export
COPY --from=base /output/pyspark_requirements.tar.gz .