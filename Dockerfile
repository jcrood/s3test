# build
FROM golang:1.18-alpine3.15 AS buildstage

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v

# runtime
FROM gcr.io/distroless/static:nonroot

COPY --from=buildstage /app/s3test /
ENTRYPOINT ["/s3test", "listObjects"]
