FROM golang:1.22-bookworm AS build

ENV GOPROXY=https://goproxy.cn,direct
ENV GOTOOLCHAIN=local
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./...
RUN go build -o /src/strata-forge .

FROM golang:1.22-bookworm
ENV GOTOOLCHAIN=local
WORKDIR /app
COPY --from=build /src/strata-forge /app/strata-forge
COPY --from=build /src/web /app/web
EXPOSE 8080
ENTRYPOINT ["/app/strata-forge"]
