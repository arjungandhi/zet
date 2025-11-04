package zet_test

import (
	"testing"

	"github.com/arjungandhi/zet/pkg/zet"
)

func TestFindNoteWithSelection(t *testing.T) {
	t.Run("empty search term", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping fzf integration test in short mode")
		}
	})

	t.Run("with search term", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping fzf integration test in short mode")
		}
	})

	t.Run("no notes", func(t *testing.T) {
		emptyNotes := []*zet.Note{}
		_, err := zet.FindNote(emptyNotes, "")
		if err == nil {
			t.Error("FindNote() with empty notes should return error")
		}
	})
}

func TestBuildFzfInput(t *testing.T) {
	notes := []*zet.Note{
		{Title: "First Note", Path: "/tmp/First Note.md", Body: "content1"},
		{Title: "Second Note", Path: "/tmp/Second Note.md", Body: "content2"},
	}

	input := zet.BuildFzfInput(notes)

	expectedNewlines := 1
	actualNewlines := 0
	for i := 0; i < len(input); i++ {
		if input[i] == '\n' {
			actualNewlines++
		}
	}

	if actualNewlines != expectedNewlines {
		t.Errorf("BuildFzfInput() produced %d newlines, want %d", actualNewlines, expectedNewlines)
	}
}

func TestParseFzfOutput(t *testing.T) {
	notes := []*zet.Note{
		{Title: "Apple Note", Path: "/tmp/Apple Note.md", Body: "content1"},
		{Title: "Banana Note", Path: "/tmp/Banana Note.md", Body: "content2"},
		{Title: "Cherry Note", Path: "/tmp/Cherry Note.md", Body: "content3"},
	}

	tests := []struct {
		name         string
		fzfOutput    string
		expectedNote *zet.Note
		wantErr      bool
	}{
		{
			name:         "valid selection - first note",
			fzfOutput:    "0\tApple Note\t/tmp/Apple Note.md",
			expectedNote: notes[0],
			wantErr:      false,
		},
		{
			name:         "valid selection - second note",
			fzfOutput:    "1\tBanana Note\t/tmp/Banana Note.md",
			expectedNote: notes[1],
			wantErr:      false,
		},
		{
			name:         "invalid index",
			fzfOutput:    "99\tNonexistent\t/tmp/none.md",
			expectedNote: nil,
			wantErr:      true,
		},
		{
			name:         "malformed output - no tabs",
			fzfOutput:    "invalid",
			expectedNote: nil,
			wantErr:      true,
		},
		{
			name:         "empty output",
			fzfOutput:    "",
			expectedNote: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := zet.ParseFzfOutput(notes, tt.fzfOutput)

			if tt.wantErr && err == nil {
				t.Error("ParseFzfOutput() expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("ParseFzfOutput() unexpected error: %v", err)
			}

			if tt.expectedNote != nil && result != tt.expectedNote {
				t.Errorf("ParseFzfOutput() = %v, want %v", result, tt.expectedNote)
			}
		})
	}
}
