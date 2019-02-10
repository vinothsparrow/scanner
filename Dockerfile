FROM golang@sha256:8dea7186cf96e6072c23bcbac842d140fe0186758bcc215acb1745f584984857 as builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
# Create appuser.
RUN go get github.com/golang/dep/cmd/dep
RUN go get -u github.com/swaggo/swag/cmd/swag
WORKDIR /go/src/github.com/vinothsparrow/scanner
ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock
# install packages
RUN dep ensure --vendor-only
COPY . .
RUN swag init
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/scanner
############################
# STEP 2 build a small image
############################
FROM iron/go
WORKDIR /go/bin/
# Copy our static executable.
COPY --from=builder /go/bin/scanner .
COPY --from=builder /go/src/github.com/vinothsparrow/scanner/docs .
EXPOSE 8000
# Run the scanner binary.
ENTRYPOINT ["./scanner"]