FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/rooommetrics ./cmd/rooommetrics

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/rooommetrics /rooommetrics
VOLUME ["/data"]
ENV ROOOM_WAL_DIR=/data
EXPOSE 9090
USER nonroot:nonroot
ENTRYPOINT ["/rooommetrics"]
