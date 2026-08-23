## Context

The README currently serves three audiences poorly mixed together: end users, contributors, and maintainers. Architecture internals (CGo, build chain) sit alongside user-facing quickstart instructions. The repo doesn't exist yet — everything must be launch-ready. See proposal.md for motivation.

Constraints:
- `docs/CONTRIBUTING.md` and `docs/MAINTAINING.md` already exist
- `internal/catalog/models.json` is the source of truth for model data
- `sample-audio/Jfk_rice_university_we_choose_to_go_to_the_moon.ogg` exists (16.4 MB, Vorbis codec, public domain)
- No Homebrew tap published; no GitHub releases exist

## Goals / Non-Goals

**Goals:**
- README rewritten for end users with SEO, quickstart, model guide, comparison
- Architecture internals extracted to `docs/ARCHITECTURE.md`
- Self-contained terminal demo via vhs + committed `.opus` clip
- GitHub repo created with description and topics

**Non-Goals:**
- Changing `docs/CONTRIBUTING.md` or `docs/MAINTAINING.md` content
- Adding new features or changing the binary
- Publishing Homebrew tap or GitHub release artifacts
- Translating README to other languages

## Decisions

### README section order
```
H1 + tagline → Demo GIF → Quickstart → What is Opuscriber? → Why Opuscriber?
→ Features → Model Guide → Usage → Installation → Documentation → License
```
**Rationale**: Quickstart and demo GIF must appear above the fold (in the first ~600px of the GitHub page). The demo GIF is a visual hook. Explanatory text comes after.

### Demo audio: transcode once, ship as .opus
The Wikimedia file is Ogg Vorbis, not Opus. Opuscriber's `pion/opus` decoder only handles Opus codec. We transcode a 30s clip once via ffmpeg to `demo/jfk-moon-speech.opus` (16kHz mono Opus), then commit it. The original stays in `sample-audio/` for provenance.
**Alternatives considered**: Adding Vorbis support to the decode pipeline — rejected as scope creep for v1 launch. The demo shows what the tool actually does.

### Model guide table: real catalog data, markdown table
Use exact `id` values from `models.json`. Disk sizes are `size_bytes` converted to human-readable. RAM is approximate (not in catalog — inferred from whisper.cpp documentation). `recommended: true` on `medium` → star marker in table.
**Alternatives considered**: Just linking to whisper.cpp model docs — rejected because it's user friction on a critical decision.

### Comparison section: prose, not table
Prose avoids the "this vs. that" combativeness of a feature table and naturally contains comparison keywords for SEO. Addresses three alternatives: raw whisper.cpp CLI, cloud APIs, GUI apps.
**Alternatives considered**: Feature matrix table — rejected for tone reasons on a pre-release project.

### Architecture: new file, simplified diagram in README
Move full architecture to `docs/ARCHITECTURE.md`. README keeps a 3-line ASCII diagram: `Opus file → decode → Whisper → TXT + SRT`. Remove all CGo, libwhisper.a, pion/opus internals.
**Rationale**: End users don't care about build chain internals. Contributors need them — put them where contributors look.

### Pre-release install: honest ordering
1. Build from source (`make build`) — what works today
2. Homebrew — "coming soon" qualifier
3. Prebuilt binaries — "coming soon" with link to future Releases
**Rationale**: Trust is critical for a new tool. Fake "Download" sections erode credibility.

### vhs for terminal recording
Use charmbracelet/vhs (same ecosystem as Bubble Tea). `.tape` file committed to repo; `demo.gif` committed as binary. `make demo` target regenerates it.
**Rationale**: vhs produces high-quality GIF renders with proper terminal font rendering. Reproducible, scriptable.

## Risks / Trade-offs

- **Demo GIF file size**: A 25-second TUI recording at 80×28 could be 1-3 MB. Acceptable for a README if reasonable. If too large, shorten the tape or reduce frame rate.
- **Model guide RAM estimates**: Not in catalog — inferred. Document them as approximate.
- **JFK clip selection**: The 30s clip must contain recognizable speech to produce interesting transcription output. The famous line starts ~7:30 into the speech.
- **Single-platform demo**: vhs runs on macOS/Linux. The demo GIF is committed, so this only matters for regenerating it.