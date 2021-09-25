FROM golang:alpine as build
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

FROM alpine
RUN mkdir /app
COPY --from=build /app/main /app/main
CMD ["/app/main"]