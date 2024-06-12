# Stage 1: Build the Go plugin
FROM heroiclabs/nakama-pluginbuilder:3.22.0 AS go-builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /backend

COPY go.mod .
COPY go.sum .
RUN go mod tidy
RUN go mod vendor
COPY ../go-server-nakima .

RUN go build --trimpath --mod=vendor --buildmode=plugin -o ./backend.so

# Stage 2: Add the plugin to Nakama
FROM registry.heroiclabs.com/heroiclabs/nakama:3.22.0

COPY --from=go-builder /backend/backend.so /nakama/data/modules/
COPY local.yml /nakama/data/

# Stage 3: Create folder and move file with data
RUN mkdir -p /nakama/data/files/core
COPY 1.0.0.json /nakama/data/files/core/1.0.0.json
ENV FILE_BASE_PATH /nakama/data/files
