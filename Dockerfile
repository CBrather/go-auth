####################
# BUILD IMAGE      #
####################
FROM alpine:latest as builder
WORKDIR /app

RUN apk update && \
  apk add --no-cache ca-certificates && \
  mkdir /user && \
  echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
  echo 'nobody:x:65534:' > /user/group

####################
# PRODUCTION IMAGE #
####################
FROM scratch AS final

EXPOSE 8080

USER nobody:nobody

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./bin/main .

ENTRYPOINT ["/main"]
