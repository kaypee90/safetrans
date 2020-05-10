# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-safetrans"
LABEL REPO="https://github.com/kaypee90/safetrans"

ENV PROJPATH=/go/src/github.com/kaypee90/safetrans

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/kaypee90/safetrans
WORKDIR /go/src/github.com/kaypee90/safetrans

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/kaypee90/safetrans"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/safetrans/bin

WORKDIR /opt/safetrans/bin

COPY --from=build-stage /go/src/github.com/kaypee90/safetrans/bin/safetrans /opt/safetrans/bin/
RUN chmod +x /opt/safetrans/bin/safetrans

# Create appuser
RUN adduser -D -g '' safetrans
USER safetrans

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/safetrans/bin/safetrans"]
