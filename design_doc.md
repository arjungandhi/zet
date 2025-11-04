# Zet Command Design Document

### Goals
- Use human-readable filenames based on note titles instead of timestamps
- Maintain simplicity with flat file structure
- Prevent duplicate notes by treating duplicate creates as edits

### File System Structure Changes

#### Structure
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

