# Copilot Instructions for magicnumber

## Project Overview
A Go library that identifies file types by reading byte signatures (magic numbers) rather than relying on file extensions or MIME types. Supports 90+ file format signatures across images, video, audio, archives, executables, documents, and text formats.

## Build, Test, and Lint Commands

### Using Task (preferred)
```bash
# List all available tasks
task --list-all

# Run tests
task test                  # Standard test run (no caching)
task testr                 # Run tests with race detection (slower)

# Lint and format
task lint                  # Run gofumpt formatter and golangci-lint

# Static analysis
task nil                   # Run nilaway static analysis for nil dereferences

# Dependencies
task update                # Update all dependencies to latest version
task patch                 # Update only patch versions of dependencies

# Documentation
task doc                   # Generate and browse module documentation locally
```

### Direct Go commands
```bash
# Test a single package
go test -count 1 -v ./archive.go ./archive_test.go

# Test a single function
go test -run TestArchive -v -count 1 ./...

# Run with verbose output and detailed test names
go test -v -count 1 ./...

# Lint check
golangci-lint run -c .golangci.yaml
```

## High-Level Architecture

### Core Design Pattern
The library uses a **matcher pattern** where each file signature is associated with a `Matcher` function that performs byte-level validation:
- `Matcher` is a function type: `func(io.ReaderAt) bool` that checks if a reader contains a specific file signature
- `Finder` is a map of `Signature -> Matcher` that contains all detection logic
- Signatures are defined as an `iota` enum starting from `ZeroByte` (-2)

### Main Components

**magicnumber.go** (core API)
- `Signature` enum: ~90 file type constants
- `Find(io.ReaderAt) Signature`: Main entry point, returns identified file type
- `MatchExt(filename, reader)`: Validates if file content matches its extension
- `New() *Finder`: Returns map of all signature matchers
- `Ext()`: Returns map of signatures to file extensions
- Helper types: `Extension`, `Finder`, `Matcher`

**Format-specific modules** (grouped by category):
- `executable.go`: DOS/Windows executables, self-extracting archives (PKLITE, PKSFX)
- `archive.go`: ZIP variants, RAR, TAR, 7z, GZip, etc. (uses PKWARE detection logic)
- `media.go`: Images (JPEG, PNG, BMP, TIFF), video (MP4, AVI, MOV), audio (MP3, WAV, FLAC, OGG)
- `cdimage.go`: CD/DVD ISO formats (ISO 9660, Nero, PowerISO, Alcohol 120)
- `text.go`: Text and document formats (UTF-8/16/32, ANSI, PDF, RTF)
- `id3.go`: ID3 tag parsing for MP3 metadata
- `synthesismusic.go`: Tracker music formats (MOD, IT, XM, MTM)

**knowns.go**: Currently empty; placeholder for additional data structures

### Detection Strategy
1. `Find()` iterates through all matchers in `Finder` map
2. For text/special cases (ANSI, plain text, XBIN), it applies secondary checks in specific order
3. Returns `Unknown` if no signature matches
4. Returns `ZeroByte` for empty files

### Testing
- Each module has a corresponding `*_test.go` file
- Tests use actual file samples from `testdata/` directory
- Pattern: `TestSignatureName` function names (e.g., `TestArchive`, `TestMSExe`)
- Use `go test -count 1` to disable caching (important for file I/O tests)

## Performance Considerations

### Key Bottleneck: Repeated `New()` Calls
**Current Issue**: Every function that detects file types creates a new `Finder` map with all matchers:
- `Find()` calls `New()` once per file
- `MatchExt()` calls `New()` once per file
- Category functions (`Archive()`, `Image()`, `Video()`, `Document()`, etc.) in knowns.go call `New()` once each

**Impact**: Building the ~90-entry matcher map and dereference with `*New()` happens on every detection call. For high-volume batch processing, this is inefficient.

### Optimization Opportunities
1. **Cache `New()` result** - Create Finder once at package init or module load, then reuse
   ```go
   var (
       defaultFinder *Finder
   )
   
   func init() {
       defaultFinder = New()
   }
   ```

2. **Lazy initialization** - If caching entire map isn't desired, lazy-initialize on first use

3. **Specialized matchers** - Category functions already filter to smaller lists, which is good for targeted detection

### Current Strength: Minimal Byte Reading
- Matchers only read necessary bytes (often 2-6 bytes from specific offsets)
- No full file buffering
- Good for large files or streaming scenarios

### Testing Impact
Tests use `go test -count 1` to avoid caching, ensuring fresh reads. This is important and should be preserved.

## Key Conventions

### Code Organization
- **One file per format category**, not one per signature type
- Matcher functions placed near related helpers (e.g., `Pklite()` near other DOS executables)
- Internal helper functions in same file as their public matchers (e.g., `NotASCII()` in text.go)

### Naming Conventions
- Matcher functions: PascalCase, match signature constant names (e.g., `MSExe()` for `MicrosoftExecutable`)
- Signature constants: Descriptive PascalCase (e.g., `PKWAREZip`, `MicrosoftExecutable`)
- `.String()` returns lowercase/hyphenated descriptive names for UX
- `.Title()` returns full format names (e.g., "JPEG File Interchange Format")

### Byte Reading Pattern
All matchers follow a consistent pattern for safety:
```go
func MatcherName(r io.ReaderAt) bool {
    const size = N      // number of bytes to read
    const offset = M    // where to read from (usually 0)
    p := make([]byte, size)
    sr := io.NewSectionReader(r, offset, size)
    if n, err := sr.Read(p); err != nil || n < size {
        return false    // handle read errors gracefully
    }
    // bytes.Equal() or bytes.Compare() or bitmask checks
    return /* validation logic */
}
```

### Testing Files
- Test files use sample files from `testdata/` directory
- Test pattern: Read file → pass to function → assert signature
- Some matchers have multiple test cases for different variants

### Linting Configuration
- Uses golangci-lint with custom config in `.golangci.yaml`
- gofumpt formatter enabled for consistent formatting
- Max complexity: 17 (cyclop), max function length: 60 lines
- Several linters disabled: exhaustive, ireturn, varnamelen, wsl variants
- Test files get leniency on complexity/function length rules

### Dependencies
- `github.com/nalgeon/be`: Likely used for byte order operations (check usage)
- `golang.org/x/text`: Text encoding support (UTF-8/16/32 detection)
- `go.uber.org/nilaway`: Static nil pointer analysis tool

### Documentation Sources
Magic number byte values sourced from:
- Gary Kessler's File Signatures Table
- Just Solve the File Format Problem (Archive Team)
- OSDev Wiki
- Wikipedia List of File Signatures
