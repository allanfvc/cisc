FROM golang:alpine as build
RUN mkdir -p /app/build
ADD . /app/
WORKDIR /app
RUN go build -o /build/cisc .

FROM alpine
RUN mkdir /app
COPY --from=build /app/build/cisc /app/cisc
CMD ["/app/cisc", "checker"]