FROM golang:1.21.4 AS builder
LABEL maintainer="Patrick Hermann patrick.hermann@sva.de"

ARG VERSION=""
ARG BUILD_DATE=""
ARG COMMIT=""
ARG GIT_PAT=""
ARG MODULE="github.com/stuttgart-things/stageTime-server"

WORKDIR /src/
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -buildvcs=false -o /bin/stageTime-server \
    -ldflags="-X ${MODULE}/internal.version=v${VERSION} -X ${MODULE}/internal.date=${BUILD_DATE} -X ${MODULE}/internal.commit=${COMMIT}"

RUN CGO_ENABLED=0 go build -o /bin/grpcCall tests/grpcCall.go

FROM alpine:3.18.4
COPY --from=builder /bin/stageTime-server /bin/stageTime-server

# FOR SERVICE TESTING
COPY --from=builder /bin/grpcCall /bin/grpcCall
COPY --from=builder /src/tests/prs.json /tmp/prs.json

ENTRYPOINT ["stageTime-server"]
