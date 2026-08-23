# ── Stage 1: Build whisper.cpp from source ──────────────────────────
FROM ubuntu:22.04 AS whisper-builder

RUN apt-get update && apt-get install -y \
    cmake build-essential git libomp-dev \
    && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/ggml-org/whisper.cpp.git /src
WORKDIR /src
RUN cmake -B build -DCMAKE_BUILD_TYPE=Release && \
    cmake --build build -j$(nproc) --target whisper-cli

# ── Stage 2: Build Go binary ────────────────────────────────────────
FROM golang:1.27-bookworm AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/opuscriber .

# ── Stage 3: Runtime image ──────────────────────────────────────────
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ffmpeg curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# whisper-cli + shared libraries
COPY --from=whisper-builder /src/build/bin/whisper-cli /usr/local/bin/
COPY --from=whisper-builder /src/build/bin/libwhisper.so.1 /usr/local/lib/
COPY --from=whisper-builder /src/build/bin/libggml.so.0 /usr/local/lib/
COPY --from=whisper-builder /src/build/bin/libggml-base.so.0 /usr/local/lib/
COPY --from=whisper-builder /src/build/bin/libggml-cpu.so.0 /usr/local/lib/
RUN ldconfig

# model download script
COPY --from=whisper-builder /src/models/download-ggml-model.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/download-ggml-model.sh

# Go binary
COPY --from=go-builder /app/opuscriber /usr/local/bin/

# Default volumes
VOLUME ["/audio/in", "/audio/out", "/models"]

ENTRYPOINT ["opuscriber"]
CMD ["--help"]
