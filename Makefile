BINARY=opuscriber
WHISPER_VERSION=v1.9.3
WHISPER_DIR=whisper.cpp
MODEL?=medium
LANG?=pt

.PHONY: build test run run-interactive clean

# ── Build ────────────────────────────────────────────────────────────────────

libwhisper.a:
	@if [ ! -d $(WHISPER_DIR) ]; then \
		echo "Cloning whisper.cpp $(WHISPER_VERSION)..."; \
		git clone --depth 1 --branch $(WHISPER_VERSION) \
			https://github.com/ggml-org/whisper.cpp.git $(WHISPER_DIR); \
	fi
	@if [ ! -f libwhisper.a ]; then \
		echo "Building whisper.cpp..."; \
		cmake -S $(WHISPER_DIR) -B $(WHISPER_DIR)/build -DBUILD_SHARED_LIBS=OFF -DCMAKE_BUILD_TYPE=Release; \
		cmake --build $(WHISPER_DIR)/build -j$$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4); \
		find $(WHISPER_DIR)/build -name 'lib*.a' -exec cp {} . \; ; \
		cp $(WHISPER_DIR)/include/whisper.h .; \
		cp $(WHISPER_DIR)/ggml/include/*.h .; \
	fi

build: libwhisper.a
	CGO_CFLAGS="-I$$(pwd)" CGO_LDFLAGS="-L$$(pwd)" CGO_ENABLED=1 go build -o $(BINARY) .

# ── Run (for development) ────────────────────────────────────────────────────

run: build
	./$(BINARY) --input ./in --output ./out --models ./models --lang $(LANG) --model $(MODEL)

run-interactive: build
	./$(BINARY) --interactive

# ── Test ─────────────────────────────────────────────────────────────────────

test:
	go test ./...

# ── Clean ────────────────────────────────────────────────────────────────────

clean:
	rm -f $(BINARY) lib*.a *.h
	rm -rf $(WHISPER_DIR)

# ── Demo ────────────────────────────────────────────────────────────────────

demo:
	vhs demo.tape