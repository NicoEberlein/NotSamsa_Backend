FROM golang:1.23.4-alpine AS build

WORKDIR /app

ARG OUTPUT_FILE
ARG GO_FILE

RUN echo ${GO_FILE} ${OUTPUT_FILE}

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/${OUTPUT_FILE} ${GO_FILE}
RUN ls -la /app/bin

FROM alpine:edge

ARG OUTPUT_FILE

RUN echo ${GO_FILE} ${OUTPUT_FILE}

COPY --from=build /app/bin/${OUTPUT_FILE} /app/bin/${OUTPUT_FILE}
RUN apk --no-cache add ca-certificates tzdata

RUN ls -la /app/bin
ENV APP_BINARY=${OUTPUT_FILE}

ENTRYPOINT "/app/bin/${APP_BINARY}"
