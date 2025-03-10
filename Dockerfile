FROM golang:1.24 AS build
WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o minalytics cmd/server/main.go

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /app/minalytics minalytics

CMD ["./server"]
