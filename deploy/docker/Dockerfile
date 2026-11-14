# syntax=docker/dockerfile:1.7

# One Dockerfile, one binary per build: pass --build-arg SERVICE=api|executor|registry|migrate.
FROM golang:1.27-alpine AS build

ARG SERVICE=api
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w \
      -X github.com/aminbenmansour/terrakube-go/internal/version.Version=${VERSION} \
      -X github.com/aminbenmansour/terrakube-go/internal/version.Commit=${COMMIT} \
      -X github.com/aminbenmansour/terrakube-go/internal/version.Date=${DATE}" \
    -o /out/terrakube ./cmd/${SERVICE}

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/terrakube /usr/local/bin/terrakube

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/terrakube"]
