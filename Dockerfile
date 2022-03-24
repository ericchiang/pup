FROM alpine:3.10 AS builder
LABEL builder=true

ENV CGO_ENABLED=0
ENV GOPATH /go

RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc
RUN go get github.com/ericchiang/pup \
 && cd $GOPATH/src/github.com/ericchiang/pup \
 && go build \
    -a \
    -ldflags '-s -w -extldflags "-static"' \
    -o /bin/pup
RUN adduser -DH user

FROM scratch

ENTRYPOINT [ "/pup" ]
CMD [ "--help" ]

COPY --from=builder /etc/passwd /etc/passwd
USER user

COPY --from=builder /bin/pup /pup
