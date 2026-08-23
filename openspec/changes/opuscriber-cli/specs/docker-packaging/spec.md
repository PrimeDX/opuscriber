## Purpose

Package the opuscriber CLI, ffmpeg, and whisper-cli into a multi-architecture Docker image (linux/amd64, linux/arm64) with model persistence via mounted volumes.

## ADDED Requirements

### Requirement: Multi-arch image
The Docker image SHALL support both linux/amd64 and linux/arm64 platforms, built via docker buildx. whisper-cli SHALL be compiled from source in the build stage for the target architecture.

#### Scenario: Build for both platforms
- **WHEN** docker buildx build --platform linux/amd64,linux/arm64 is run
- **THEN** the image SHALL build successfully for both architectures

#### Scenario: Native ARM performance
- **WHEN** the image runs on linux/arm64 hardware
- **THEN** whisper-cli SHALL be a native ARM64 binary (not emulated x86)

### Requirement: Volume-mounted model
The whisper model file SHALL be downloaded once and persisted in a models/ directory mounted as a Docker volume.

#### Scenario: Model path
- **WHEN** the container runs with -v $(pwd)/models:/models
- **THEN** the system SHALL look for the model at /models/ggml-medium.bin

#### Scenario: Model download
- **WHEN** the user runs the download-ggml-model.sh script in the image
- **THEN** the model SHALL be written to the mounted /models/ directory

### Requirement: Volume-mounted input and output
Input audio files SHALL be mounted at /audio/in and output transcripts at /audio/out, both as Docker volumes.

#### Scenario: Input mount
- **WHEN** the container runs with -v $(pwd)/in:/audio/in
- **THEN** all .opus/.ogg/.oga files in $(pwd)/in SHALL be processed

#### Scenario: Output mount
- **WHEN** transcription completes
- **THEN** .txt and .srt files SHALL appear in the mounted /audio/out directory

### Requirement: Minimal runtime image
The final runtime image SHALL contain only: Ubuntu 22.04, ffmpeg, curl, ca-certificates, whisper-cli binary with its shared libraries, the download-ggml-model.sh script, and the Go binary.

#### Scenario: Multi-stage build
- **WHEN** the Dockerfile builds
- **THEN** whisper.cpp SHALL compile in a builder stage and only the compiled artifacts SHALL be copied to the runtime stage

### Requirement: Model download script bundled
The download-ggml-model.sh script from whisper.cpp SHALL be included in the runtime image at /usr/local/bin/.

#### Scenario: Script available
- **WHEN** the user runs docker run opuscriber sh -c "download-ggml-model.sh medium /models"
- **THEN** the script SHALL be present and executable
