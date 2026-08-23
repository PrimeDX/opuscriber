## Why

The README describes Opuscriber from the inside out (what it *is*) rather than from the outside in (what problems people search for). Key SEO terms — "speech-to-text", "audio transcription", "whisper", "offline STT" — appear zero times. The tool is presented as a WhatsApp/Telegram voice note transcriber, when it's actually a general-purpose offline Opus audio transcriber. Architecture internals clutter user-facing content. The repo doesn't exist yet — everything needs to be launch-ready.

## What Changes

- **Rewrite README.md**: SEO-optimized opening, reordered sections, new "Why Opuscriber?" comparison prose, model guide table from real catalog data, honest pre-release install instructions
- **Extract architecture to `docs/ARCHITECTURE.md`**: Move CGo, whisper.cpp build chain, pion/opus internals out of README. Keep a simplified 3-step pipeline diagram.
- **Add terminal demo**: `demo.tape` (vhs script) + committed `demo.gif` showing TUI transcribing JFK moon speech
- **Add sample audio**: `sample-audio/` already contains the Wikimedia source; add a 30s Opus clip at `demo/jfk-moon-speech.opus` for the demo to transcribe
- **Create GitHub repo**: `primedx/opuscriber` with description and topics for GitHub search discoverability
- **`make demo`**: Optional Makefile target to record/update the demo GIF

## Capabilities

### New Capabilities
- `readme-content`: README.md narrative, structure, SEO, and model guide
- `repo-launch`: GitHub repo creation, description, topics, demo GIF