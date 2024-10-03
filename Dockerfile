# syntax=docker/dockerfile:1

FROM golang:1.22 as build
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o simpler-tha ./cmd/main.go

FROM scratch
COPY --from=build /app/simpler-tha .
COPY --from=build /app/app.env .
CMD ["./simpler-tha"]
