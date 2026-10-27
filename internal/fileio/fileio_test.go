package fileio_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

func copyAndChecksumFile(src string) error {
	dir, err := fileio.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	dest := filepath.Join(dir, "dest.txt")
	err = fileio.CopyFile(src, dest)
	if err != nil {
		return err
	}

	equals, err := fileio.CompareFiles(src, dest)
	if err != nil {
		return err
	}
	if equals == false {
		return fmt.Errorf("destChecksum: %s != srcChecksum: %s", dest, src)
	}
	return nil
}

func TestCopyFile(t *testing.T) {
	err := copyAndChecksumFile(filepath.Join("testdata", "smallexe"))
	if err != nil {
		t.Fatal(err)
	}
	err = copyAndChecksumFile(filepath.Join("testdata", "TempleOS.ISO"))
	if err != nil {
		t.Fatal(err)
	}
	err = copyAndChecksumFile(filepath.Join("testdata", "empty.txt"))
	if err != nil {
		t.Fatal(err)
	}

	err = fileio.CopyFile(filepath.Join("testdata", "smallexe"), filepath.Join("pathdoesnotexists", "pathdoesnotexists"))
	if err == nil {
		t.Fatal("CopyFile should not work if path does not exists")
	}
	err = fileio.CopyFile(filepath.Join("pathdoesnotexists", "pathdoesnotexists"), filepath.Join("testdata", "smallexe"))
	if err == nil {
		t.Fatal("CopyFile should not work if path does not exists")
	}

}

func verifyChecksumFile(src, expectedChecksum string) error {
	checksum, err := fileio.ChecksumFile(src)
	if err != nil {
		return err
	}
	if checksum != expectedChecksum {
		return fmt.Errorf("checksum: %s != expectedChecksum: %s", checksum, expectedChecksum)
	}
	return nil
}

func TestChecksum(t *testing.T) {
	err := verifyChecksumFile(filepath.Join("testdata", "empty.txt"),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Fatal(err)
	}
	err = verifyChecksumFile(filepath.Join("testdata", "TempleOS.ISO"),
		"5d0fc944e5d89c155c0fc17c148646715bc1db6fa5750c0b913772cfec19ba26")
	if err != nil {
		t.Fatal(err)
	}

	// Test wrong checksum:
	err = verifyChecksumFile(filepath.Join("testdata", "TempleOS.ISO"),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err == nil {
		t.Fatal(fmt.Errorf("verifyChecksumFile returned not nil on a bad checksum"))
	}
}

func TestCompareFileChecksum(t *testing.T) {
	fileA := filepath.Join("testdata", "TempleOS.ISO")
	fileB := filepath.Join("testdata", "smallexe")
	fileC := filepath.Join("testdata", "doesNotExists")

	equals, err := fileio.CompareFiles(fileA, fileA)
	if err != nil {
		t.Fatal(err)
	}
	if equals == false {
		t.Fatal("Should be equal")
	}

	equals, err = fileio.CompareFiles(fileA, fileB)
	if err != nil {
		t.Fatal(err)
	}
	if equals == true {
		t.Fatal("Should be unequal")
	}

	_, err = fileio.CompareFiles(fileA, fileC)
	if err == nil {
		t.Fatal("fileio.CompareFiles(fileA, fileC) should return an error")
	}
	_, err = fileio.CompareFiles(fileC, fileB)
	if err == nil {
		t.Fatal("fileio.CompareFiles(fileC, fileB) should return an error")
	}
}

func TestFileExists(t *testing.T) {
	fileA := filepath.Join("testdata", "TempleOS.ISO")
	fileB := filepath.Join("testdata", "doesNotExists")

	if !fileio.FileExists(fileA) {
		t.Error("fileio.FileExists(fileA) exists")
	}
	if fileio.FileExists(fileB) {
		t.Error("fileio.FileExists(fileA) does not exists")
	}
}

func TestSafeJoinEscape(t *testing.T) {
	dest := filepath.Clean(os.TempDir())

	for _, name := range []string{"file.txt", "sub/dir/file.txt", "sub/../file.txt"} {
		got, err := fileio.SafeJoin(dest, name)
		if err != nil {
			t.Errorf("SafeJoin(%q) returned error: %v", name, err)
			continue
		}
		if got != dest && !strings.HasPrefix(got, dest+string(os.PathSeparator)) {
			t.Errorf("SafeJoin(%q) = %q, which is outside dest %q", name, got, dest)
		}
	}

	for _, name := range []string{"../escape.txt", "../../escape.txt", "sub/../../escape.txt"} {
		if _, err := fileio.SafeJoin(dest, name); err == nil {
			t.Errorf("SafeJoin(%q) returned nil error, want path-traversal rejection", name)
		}
	}
}

func TestSafeJoin(t *testing.T) {
	dest := filepath.Join("root", "uploads")

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple filename",
			input: "file.txt",
			want:  filepath.Join(dest, "file.txt"),
		},
		{
			name:  "nested filename",
			input: filepath.Join("images", "photo.png"),
			want:  filepath.Join(dest, "images", "photo.png"),
		},
		{
			name:  "multiple nested directories",
			input: filepath.Join("a", "b", "c", "file.txt"),
			want:  filepath.Join(dest, "a", "b", "c", "file.txt"),
		},
		{
			name:  "current directory component",
			input: filepath.Join(".", "file.txt"),
			want:  filepath.Join(dest, "file.txt"),
		},
		{
			name:  "internal parent component",
			input: filepath.Join("sub", "..", "file.txt"),
			want:  filepath.Join(dest, "file.txt"),
		},
		{
			name:  "empty filename",
			input: "",
			want:  dest,
		},
		{
			name:  "destination itself",
			input: ".",
			want:  dest,
		},
		{
			name:    "parent traversal",
			input:   "..",
			wantErr: true,
		},
		{
			name:    "parent traversal with filename",
			input:   filepath.Join("..", "secret.txt"),
			wantErr: true,
		},
		{
			name:    "multiple parent traversal",
			input:   filepath.Join("..", "..", "secret.txt"),
			wantErr: true,
		},
		{
			name:    "traversal after nested directory",
			input:   filepath.Join("sub", "..", "..", "secret.txt"),
			wantErr: true,
		},
		{
			name:    "directory prefix lookalike",
			input:   filepath.Join("..", "uploads-evil", "file.txt"),
			wantErr: true,
		},
		{
			name:  "nested directory with redundant components",
			input: filepath.Join("sub", ".", "nested", "file.txt"),
			want:  filepath.Join(dest, "sub", "nested", "file.txt"),
		},
		{
			name:  "nested directory with internal traversal",
			input: filepath.Join("a", "b", "..", "file.txt"),
			want:  filepath.Join(dest, "a", "file.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileio.SafeJoin(dest, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("SafeJoin(%q, %q) expected error, got %q",
						dest, tt.input, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("SafeJoin(%q, %q) unexpected error: %v",
					dest, tt.input, err)
			}

			if got != tt.want {
				t.Errorf("SafeJoin(%q, %q) = %q, want %q",
					dest, tt.input, got, tt.want)
			}
		})
	}
}

func TestSafeJoinCleansDestination(t *testing.T) {
	dest := filepath.Join("root", "uploads", "..", "uploads")
	name := "file.txt"

	got, err := fileio.SafeJoin(dest, name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := filepath.Join("root", "uploads", "file.txt")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSafeJoinRejectsTraversalAtDifferentDepths(t *testing.T) {
	dest := filepath.Join("root", "uploads")

	tests := []struct {
		depth int
	}{
		{depth: 1},
		{depth: 2},
		{depth: 3},
		{depth: 5},
		{depth: 10},
	}

	for _, tt := range tests {
		t.Run(strings.Repeat("parent_", tt.depth), func(t *testing.T) {
			parts := make([]string, tt.depth)
			for i := range parts {
				parts[i] = ".."
			}
			parts = append(parts, "secret.txt")

			name := filepath.Join(parts...)
			got, err := fileio.SafeJoin(dest, name)

			if err == nil {
				t.Errorf("expected traversal to be rejected, got %q", got)
			}
		})
	}
}
