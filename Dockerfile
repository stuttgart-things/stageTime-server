FROM eu.gcr.io/stuttgart-things/sthings-golang:1.22 AS builder
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

FROM alpine:3.17.0
COPY --from=builder /bin/stageTime-server /bin/stageTime-server

ENTRYPOINT ["stageTime-server"]
