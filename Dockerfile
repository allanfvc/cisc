FROM golang:alpine as build
RUN mkdir -p /app/build
ADD . /app/
WORKDIR /app
RUN go build -o /build/cisc cmd/main.go

FROM alpine
RUN mkdir /app
COPY --from=build /app/build/cisc /app/cisc
RUN chmod 755 /app/cisc
CMD ["/app/cisc", "checker"]