# MAS Testing & Coverage Guide

This document explains how to run tests and generate coverage reports for the MAS (Multi-Agent System) project, and highlights important gotchas for Go users.

---

## Running Tests

From the MAS project root:

```bash
cd agent
# Run all tests in the agent package
go test -v
```

Or, from the project root:

```bash
go test -v ./agent
```

---

## Generating Coverage Reports

**Important:** Never place coverage output files (e.g., `*.out`) in your Go module root or any package directory. Go may interpret these as packages, causing cryptic errors like:

```
# .out
no required module provides package .out; to add it:
        go get .out
FAIL    .out [setup failed]
```

### Recommended Coverage Workflow

1. **Create a temp or build directory for coverage output:**
   - On Windows: `C:/Temp` or a dedicated `build/` directory
   - On Unix: `/tmp` or `build/`

2. **Run coverage:**
   ```bash
   go test -coverprofile=C:/Temp/mas_coverage.out ./agent
   # Or
   go test -coverprofile=build/mas_coverage.out ./agent
   ```

3. **View coverage summary:**
   ```bash
   go tool cover -func C:/Temp/mas_coverage.out
   # Or
   go tool cover -func build/mas_coverage.out
   ```

---

## Troubleshooting

### If you see the ".out" error:
- Delete any `.out` files or directories in your project tree.
- Move your coverage or build output to a temp or build directory outside the Go package tree.
- Rerun your test or coverage command.

---

## Example `.gitignore` for Build Artifacts

```
# Build and coverage artifacts
/build/
*.out
*.exe
*.test
*.log
```

---

## Additional Notes

- Never build Go binaries directly to your module root.
- Keep all build and coverage artifacts in a dedicated directory (e.g., `build/`, `coverage/`, `C:/Temp`).
- This will prevent Go from misinterpreting artifacts as packages and keep your workflow clean.

---

For more details, see the [README](README.md) or contact the MAS maintainers.
