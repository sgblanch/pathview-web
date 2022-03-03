FROM node:lts-alpine AS node

WORKDIR /app/resources

COPY resources/package.json resources/package-lock.json ./
RUN npm install \
    && install -m 0755 -d /app/static

COPY resources/ ./
RUN  npm run build

########################################################################
FROM golang:alpine AS builder

WORKDIR /app

RUN apk update \
    && apk add --no-cache --virtual .build-deps ca-certificates tzdata \
    && apk add --no-cache git \
    && update-ca-certificates \
    && adduser -D -g '' pathview \
    && install -m 1777 -o root -g root -d /build/tmp \
    && install -m 0755 -o root -g root -d /build/app /build/etc/ssl/certs /build/usr/share \
    && install -m 0755 -o pathview -g pathview -d /build/opt/kegg \
    && tar c -C / etc/passwd etc/group etc/ssl/certs/ca-certificates.crt usr/share/zoneinfo | tar x -C /build  \
    && apk del .build-deps \
    && rm -rf /var/cache/apk/*

COPY go.mod go.sum ./
RUN go mod download

COPY --from=node /app/static/ /build/static/
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /build/app/ \
    && tar c static template | tar x -C /build/app \
    && rm -rf /go/pkg/* /go/src/*

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
