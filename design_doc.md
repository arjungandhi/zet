# Zet Command Design Document

## Current Design

### Overview
Zet is a CLI tool for managing a personal Zettelkasten (knowledge management system) with interconnected markdown notes. Built on the Bonzai framework for composable command structures.

### Architecture

#### Directory Structure
```
zet/
├── cmd/zet/main.go           # Entry point (12 lines)
├── pkg/zet/
│   ├── cmd.go                # Command definitions & CLI logic (235 lines)
│   ├── zet.go                # Core business logic (115 lines)
│   └── zet_test.go           # Unit tests (104 lines)
├── go.mod                    # Dependencies
└── magefile.go               # Build automation
```

#### Data Model
```go
type Note struct {
    Title string  // Extracted from "# Title" markdown heading
    Path  string  // Full path to README.md file
    Body  string  // Complete file contents
}
```

#### File System Structure
```
$ZETDIR/
├── 20231103T211524Z/         # ISO second timestamp
│   └── README.md             # Note content with "# Title" heading
├── 20231103T212015Z/
│   └── README.md
└── ...
```

### Commands

#### 1. Default Command (Interactive Open)
**Usage**: `zet [search_term]`
**Location**: `pkg/zet/cmd.go:112-135`
**Flow**:
1. Get ZETDIR from environment
2. List all notes
3. Use fzf to select note (with optional search term)
4. Open selected note in $EDITOR

#### 2. List Command
**Usage**: `zet list`
**Location**: `pkg/zet/cmd.go:72-90`
**Flow**:
1. Get ZETDIR from environment
2. List all notes
3. Print each note title to stdout

#### 3. New Command
**Usage**: `zet new [title]`
**Location**: `pkg/zet/cmd.go:92-110`
**Flow**:
1. Get ZETDIR from environment
2. Join all args as title
3. Create note with ISO second timestamp directory
4. Create README.md with "# Title" heading
5. Open in editor

#### 4. Delete Command
**Usage**: `zet delete [search_term]`
**Location**: `pkg/zet/cmd.go:137-159`
**Flow**:
1. Get ZETDIR from environment
2. List all notes
3. Use fzf to select note
4. Delete entire note directory
5. Print confirmation

#### 5. Render Command
**Usage**: `zet render [search_term]`
**Location**: `pkg/zet/cmd.go:161-179`
**Flow**:
1. Get ZETDIR from environment
2. List all notes
3. Use fzf to select note
4. Render with glow markdown viewer

### Core Functions

#### In `pkg/zet/zet.go`:
- `CreateNote(dir, title string) (*Note, error)` - Create timestamped note directory and README.md
- `DeleteNote(dir string, note *Note) error` - Remove note directory
- `ListNotes(dir string) ([]*Note, error)` - Walk directory tree and collect all notes

#### In `pkg/zet/cmd.go`:
- `getZetDir() string` - Get ZETDIR from environment or panic
- `findNote(notes []*Note, search string) *Note` - Use fzf to select note
- `editNote(note *Note)` - Open note in $EDITOR
- `renderNote(note *Note)` - Display with glow

### Dependencies
- **Bonzai** (`github.com/rwxrob/bonzai`) - Command framework
- **uniq-go** (`github.com/rwxrob/uniq-go`) - ISO second timestamp generation
- **fzf** (external binary) - Fuzzy finding
- **glow** (external binary) - Markdown rendering

### Configuration
- `ZETDIR` - Environment variable pointing to notes directory (required)
- `EDITOR` - Editor command (defaults to "vi")

### Known Issues

#### Critical (Will Not Compile)
1. **Undefined `fzf` package** (`cmd.go:202`)
   - References `fzf.ParseOptions()` but package not imported
2. **Undefined `fzf_args` variable** (`cmd.go:218`)
   - Uses undefined variable in exec.Command

#### Code Quality
3. Unused `options` variable (`cmd.go:202`)
4. Uses deprecated `io/ioutil` instead of `os` package
5. Hardcoded "README.md" filename and "# " prefix
6. No error handling in `editNote()` function
7. No bounds checking on string split (`cmd.go:226`)

#### Architecture
8. High code duplication across commands (all use same ZETDIR → ListNotes → findNote pattern)
9. Mixed concerns in `cmd.go` (CLI + business logic)
10. No configuration system (only single env var)
11. Minimal test coverage (no CLI command tests)
12. No validation that external dependencies (fzf, glow) exist
13. No logging or debug mode

---

## Proposed New Design

### Goals
- Use human-readable filenames based on note titles instead of timestamps
- Maintain simplicity with flat file structure
- Prevent duplicate notes by treating duplicate creates as edits
- Rely on filesystem timestamps instead of custom metadata

### File System Structure Changes

#### New Structure
```
$ZETDIR/
├── My First Note.md
├── Meeting Notes.md
├── Ideas and Thoughts.md
└── Project Planning.md
```

#### Filename Rules
- Use note title directly as filename with `.md` extension
- **Allow spaces** in filenames
- **Remove all special characters** (e.g., `:`, `!`, `?`, `&`, `/`, `\`, etc.)
- **Preserve case** from user input
- Examples:
  - `"My New Note"` → `My New Note.md`
  - `"My New Note: Ideas & Thoughts!"` → `My New Note Ideas Thoughts.md`
  - `"What's the plan?"` → `Whats the plan.md`

#### Behavior Changes
- **Duplicate titles**: If note already exists, open it for editing (treat `zet new` as edit)
- **Timestamps**: Use filesystem timestamps only (no frontmatter or metadata)
- **Sorting**: Alphabetical order by title in all list views

#### Data Model Changes
```go
type Note struct {
    Title string  // Title from filename (without .md extension)
    Path  string  // Full path to .md file
    Body  string  // Complete file contents
}
```

### Command Changes

#### Merged New/Edit Command
Instead of separate `new` and `edit` commands, have a single command:

**`zet new [title]`** or **`zet edit [title]`** (aliases for same behavior):
- If note with that title exists → open it for editing
- If note doesn't exist → create it and open for editing

This simplifies the UX: users don't need to remember if a note exists or not.

#### Updated Command Set
1. **`zet [search_term]`** - Interactive open (unchanged, uses fzf)
2. **`zet new [title]`** - Create or edit note by title
3. **`zet list`** - List all notes (alphabetically sorted)
4. **`zet delete [search_term]`** - Delete a note (unchanged, uses fzf)
5. **`zet render [search_term]`** - Render with glow (unchanged, uses fzf)

### Architecture Changes

#### Current Problems
1. **Mixed concerns**: `cmd.go` has CLI commands + helper functions mixed together
2. **Code duplication**: Every command repeats getZetDir → ListNotes → findNote pattern
3. **Unclear separation**: Business logic scattered between `cmd.go` and `zet.go`
4. **External dependencies**: fzf, glow, editor logic embedded in command layer

#### Proposed Clean Architecture

```
pkg/zet/
├── cmd.go           # CLI layer - Bonzai command definitions only
├── zet.go           # Business logic - note operations
├── finder.go        # Note selection logic (fzf integration)
├── config.go        # Configuration management (ZETDIR, EDITOR, etc.)
└── filesystem.go    # File operations (sanitize, read, write, list)
```

**Layer Responsibilities:**

1. **`cmd.go` (CLI Layer)**
   - Define Bonzai commands
   - Parse arguments
   - Call zet layer
   - Handle errors and output to user
   - Should be thin - just routing

2. **`zet.go` (Business Logic)**
   - `CreateOrEditNote(title)` - Unified create/edit logic
   - `DeleteNote(note)` - Delete operations
   - `ListNotes()` - Get all notes
   - `OpenNote(search)` - Find and open note
   - `RenderNote(search)` - Find and render note
   - Orchestrates finder, config, and filesystem layers

3. **`finder.go` (Selection Logic)**
   - `FindNote(notes, search)` - fzf integration
   - Handles external fzf process
   - Returns selected note

4. **`config.go` (Configuration)**
   - `GetZetDir()` - Get notes directory
   - `GetEditor()` - Get editor command
   - `GetRenderer()` - Get markdown renderer
   - Single source of truth for configuration

5. **`filesystem.go` (File Operations)**
   - `SanitizeFilename(title)` - Clean title for filesystem
   - `ReadNote(path)` - Read note from disk
   - `WriteNote(path, content)` - Write note to disk
   - `ListNoteFiles(dir)` - Scan directory for .md files
   - `NoteExists(dir, title)` - Check if note exists
   - Pure functions, no business logic

#### Benefits
- **Single Responsibility**: Each file has one clear purpose
- **Testability**: Easy to test each layer independently
- **DRY**: Common patterns extracted to service layer
- **Maintainability**: Easy to find where logic lives
- **Extensibility**: Easy to add new commands or storage backends

#### Key Changes
- Keep `zet.go` as business logic layer (refactor existing code)
- Extract fzf logic to new `finder.go`
- Extract config logic to new `config.go`
- Extract filesystem logic to new `filesystem.go`
- Simplify `cmd.go` to just command definitions
- `zet.go` orchestrates finder, config, and filesystem layers

### Migration Plan

#### For Existing Notes
Need to migrate from:
```
20231103T211524Z/README.md  (with "# My Note" inside)
```
To:
```
My Note.md  (with "# My Note" inside)
```

#### Migration Strategy Options
1. **Manual migration**: User migrates their own notes
2. **Migration command**: Add `zet migrate` command to convert old structure
3. **Automatic detection**: Detect old structure and convert on-the-fly

**Recommendation**: Implement `zet migrate` command that:
- Scans for old timestamp directories
- Extracts title from README.md heading
- Renames to new structure
- Provides dry-run option

### Implementation Steps

#### 1. Create New Files (Clean Architecture)

**`pkg/zet/filesystem.go`** - File operations layer
```go
// SanitizeFilename removes special characters from title, keeps spaces
func SanitizeFilename(title string) string {
    // Remove all non-alphanumeric except spaces
    // Trim leading/trailing spaces
    // Return: "My Note Title"
}

// NoteExists checks if a note file exists
func NoteExists(dir, title string) bool {
    sanitized := SanitizeFilename(title)
    path := filepath.Join(dir, sanitized + ".md")
    _, err := os.Stat(path)
    return err == nil
}

// ReadNote reads a note file and returns Note struct
func ReadNote(path string) (*Note, error) {
    // Extract title from filename (remove .md extension)
    // Read file contents as Body
    // Return &Note{Title: title, Path: path, Body: content}
}

// WriteNote creates/updates a note file
func WriteNote(dir, title, content string) error {
    sanitized := SanitizeFilename(title)
    path := filepath.Join(dir, sanitized + ".md")
    return os.WriteFile(path, []byte(content), 0644)
}

// ListNoteFiles scans directory for .md files, returns sorted alphabetically
func ListNoteFiles(dir string) ([]string, error) {
    // Glob for *.md files
    // Sort alphabetically
    // Return file paths
}
```

**`pkg/zet/config.go`** - Configuration layer
```go
// GetZetDir returns ZETDIR from environment or error
func GetZetDir() (string, error) {
    dir := os.Getenv("ZETDIR")
    if dir == "" {
        return "", fmt.Errorf("ZETDIR environment variable not set")
    }
    return dir, nil
}

// GetEditor returns EDITOR from environment or default
func GetEditor() string {
    editor := os.Getenv("EDITOR")
    if editor == "" {
        return "vi"
    }
    return editor
}

// GetRenderer returns markdown renderer command
func GetRenderer() string {
    return "glow" // Could make configurable later
}
```

**`pkg/zet/finder.go`** - fzf integration layer
```go
// FindNote uses fzf to interactively select a note
func FindNote(notes []*Note, searchTerm string) (*Note, error) {
    // Create fzf input with note titles
    // Run fzf with search term (if provided)
    // Parse selection and return matching Note
}
```

#### 2. Update Existing Files

**`pkg/zet/zet.go`** - Business logic (orchestration)

Update Note struct:
```go
type Note struct {
    Title string  // Now from filename, not file content
    Path  string  // Full path to .md file
    Body  string  // File contents (no title heading)
}
```

Replace `CreateNote()` with `CreateOrEditNote()`:
```go
func CreateOrEditNote(title string) error {
    dir, err := config.GetZetDir()
    if err != nil {
        return err
    }

    sanitized := filesystem.SanitizeFilename(title)
    path := filepath.Join(dir, sanitized + ".md")

    // Check if note exists
    if filesystem.NoteExists(dir, title) {
        // Edit existing note
        return openInEditor(path)
    }

    // Create new note with empty content (or just title as content?)
    err = filesystem.WriteNote(dir, title, "")
    if err != nil {
        return err
    }

    // Open in editor
    return openInEditor(path)
}
```

Update `ListNotes()`:
```go
func ListNotes() ([]*Note, error) {
    dir, err := config.GetZetDir()
    if err != nil {
        return nil, err
    }

    // Get all .md files
    files, err := filesystem.ListNoteFiles(dir)
    if err != nil {
        return nil, err
    }

    // Convert to Note structs
    var notes []*Note
    for _, path := range files {
        note, err := filesystem.ReadNote(path)
        if err != nil {
            continue // Skip errors
        }
        notes = append(notes, note)
    }

    return notes, nil
}
```

Update `DeleteNote()`:
```go
func DeleteNote(note *Note) error {
    // Simply remove the .md file (no directory)
    return os.Remove(note.Path)
}
```

Add new functions:
```go
func OpenNote(searchTerm string) error {
    notes, err := ListNotes()
    if err != nil {
        return err
    }

    note, err := finder.FindNote(notes, searchTerm)
    if err != nil {
        return err
    }

    editor := config.GetEditor()
    return openInEditor(note.Path)
}

func RenderNote(searchTerm string) error {
    notes, err := ListNotes()
    if err != nil {
        return err
    }

    note, err := finder.FindNote(notes, searchTerm)
    if err != nil {
        return err
    }

    renderer := config.GetRenderer()
    cmd := exec.Command(renderer, note.Path)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func openInEditor(path string) error {
    editor := config.GetEditor()
    // Use bonzai.SysExec or exec.Command
    return bonzai.SysExec(editor, path)
}
```

**`pkg/zet/cmd.go`** - CLI layer (simplified)

Update all commands to call zet layer:
```go
var Cmd = &bonzai.Cmd{
    Name: "zet",
    Commands: []*bonzai.Cmd{
        helpCmd,
        listCmd,
        newCmd,
        deleteCmd,
        renderCmd,
    },
    Call: defaultCmd, // Interactive open
}

var defaultCmd = func(args []string) error {
    search := strings.Join(args, " ")
    return zet.OpenNote(search)
}

var listCmd = &bonzai.Cmd{
    Name: "list",
    Call: func(args []string) error {
        notes, err := zet.ListNotes()
        if err != nil {
            return err
        }
        for _, note := range notes {
            fmt.Println(note.Title)
        }
        return nil
    },
}

var newCmd = &bonzai.Cmd{
    Name: "new",
    Call: func(args []string) error {
        title := strings.Join(args, " ")
        if title == "" {
            return fmt.Errorf("title required")
        }
        return zet.CreateOrEditNote(title)
    },
}

var deleteCmd = &bonzai.Cmd{
    Name: "delete",
    Call: func(args []string) error {
        search := strings.Join(args, " ")
        notes, err := zet.ListNotes()
        if err != nil {
            return err
        }

        note, err := finder.FindNote(notes, search)
        if err != nil {
            return err
        }

        err = zet.DeleteNote(note)
        if err != nil {
            return err
        }

        fmt.Printf("Deleted: %s\n", note.Title)
        return nil
    },
}

var renderCmd = &bonzai.Cmd{
    Name: "render",
    Call: func(args []string) error {
        search := strings.Join(args, " ")
        return zet.RenderNote(search)
    },
}
```

#### 3. Key Implementation Details

**Title Extraction:**
- Old: Parse first line, strip "# " prefix
- New: Extract from filename, remove ".md" extension
  ```go
  title := strings.TrimSuffix(filepath.Base(path), ".md")
  ```

**File Structure:**
- Old: `$ZETDIR/20231103T211524Z/README.md`
- New: `$ZETDIR/My Note Title.md`

**Sanitization Logic:**
```go
func SanitizeFilename(title string) string {
    // Use regex: [^a-zA-Z0-9 ] to remove non-alphanumeric except spaces
    re := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
    sanitized := re.ReplaceAllString(title, "")

    // Replace multiple spaces with single space
    sanitized = strings.Join(strings.Fields(sanitized), " ")

    // Trim spaces
    return strings.TrimSpace(sanitized)
}
```

**Alphabetical Sorting:**
```go
import "sort"

func ListNoteFiles(dir string) ([]string, error) {
    files, err := filepath.Glob(filepath.Join(dir, "*.md"))
    if err != nil {
        return nil, err
    }

    sort.Strings(files)
    return files, nil
}
```

#### 4. Deprecated Package Fixes
- Replace `io/ioutil.WriteFile` with `os.WriteFile`
- Replace `io/ioutil.ReadFile` with `os.ReadFile`

#### 5. Bug Fixes
- Remove undefined `fzf` package references in cmd.go
- Fix `fzf_args` undefined variable
- Remove unused `options` variable

#### 6. Migration Status
- ✅ Migration script created and executed
- ✅ 207 notes migrated from timestamp dirs to title-based files
- ✅ Empty timestamp directories removed
- ✅ Title headings removed from all note files
- Next: Implement new architecture
