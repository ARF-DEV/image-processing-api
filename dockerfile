FROM golang:1.23 AS build

COPY . . 
RUN ls -la
RUN go build -o /main 
RUN ls -la

FROM debian:bookworm-slim 
COPY --from=build /main /main
RUN apt-get update && apt-get install -y ca-certificates

EXPOSE 8080
CMD ["./main"]