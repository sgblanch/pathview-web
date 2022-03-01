FROM node:lts-alpine AS node

COPY resources/ /web/resources/
WORKDIR /web/resources

RUN npm install && \
    install -m 0755 -d /web/static && \
    npm run build && \
    rm -rf node_modules /root/.npm

########################################################################
FROM golang:alpine AS builder

WORKDIR /go/src/github.com/sgblanch/pathview-web

RUN apk update && \
    apk add --no-cache --virtual .build-deps ca-certificates tzdata && \
    apk add --no-cache git && \
    update-ca-certificates && \
    adduser -D -g '' pathview && \
    install -m 1777 -o root -g root -d /build/tmp && \
    install -m 0755 -o root -g root -d /build/app /build/etc/ssl/certs /build/usr/share && \
    install -m 0755 -o pathview -g pathview -d /build/opt/kegg && \
    tar c -C / etc/passwd etc/group etc/ssl/certs/ca-certificates.crt usr/share/zoneinfo | tar x -C /build  && \
    apk del .build-deps && \
    rm -rf /var/cache/apk/*

COPY . .
COPY --from=node /web/static/ static/

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /build/app/  && \
    tar c static template | tar x -C /build/app && \
    rm -rf /go/pkg/* /go/src/*

########################################################################
FROM scratch

LABEL description="Web front-end to pathview R package" \
    url="https://github.com/sgblanch/pathview-web" \
    maintainer="Steven Blanchard <sgblanch@uncc.edu>"

COPY --from=builder /build/ /

USER pathview

WORKDIR /app
VOLUME /opt/kegg
ENTRYPOINT ["/app/pathview-web"]
CMD ["server"]
