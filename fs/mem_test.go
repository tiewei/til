package fs_test

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"testing"

	. "bridgedl/fs"
)

func TestMemFS(t *testing.T) {
	const filePath = "/path/to/file.txt"
	const fileData = "aaabb"

	t.Run("create and read file", func(t *testing.T) {
		mfs := NewMemFS()

		if err := mfs.CreateFile(filePath, []byte(fileData)); err != nil {
			t.Fatal("Error creating file:", err)
		}

		f, err := mfs.Open(filePath)
		if err != nil {
			t.Fatal("Error opening file:", err)
		}

		const bufSize = 3
		buf := make([]byte, bufSize)

		// read first 3 bytes
		n, err := f.Read(buf)
		if err != nil {
			t.Fatal("Error reading file:", err)
		}
		if expectN := bufSize; n != expectN {
			t.Fatalf("Expected %d bytes to be read, returned %d (buffer content: %q)", expectN, n, buf)
		}
		if len(buf) != bufSize {
			t.Fatalf("Buffer size has changed from %d to %d (buffer content: %q)", bufSize, n, buf)
		}
		if expectBuf := []byte{'a', 'a', 'a'}; !bytes.Equal(buf, expectBuf) {
			t.Fatalf("Expected buffer content to be %q, got %q", expectBuf, buf)
		}

		// read next 3 bytes, but only 2 bytes remain to be read
		n, err = f.Read(buf)
		if err != nil {
			t.Fatal("Error reading file:", err)
		}
		if expectN := len(fileData) - bufSize; n != expectN {
			t.Fatalf("Expected %d bytes to be read, returned %d (buffer content: %q)", expectN, n, buf)
		}
		if len(buf) != bufSize {
			t.Fatalf("Buffer size has changed from %d to %d (buffer content: %q)", bufSize, n, buf)
		}
		if expectBuf := []byte{'b', 'b', 'a'}; !bytes.Equal(buf, expectBuf) {
			t.Fatalf("Expected buffer content to be %q, got %q", expectBuf, buf)
		}

		// no more byte to read, expect EOF
		n, err = f.Read(buf)
		if n != 0 {
			t.Fatalf("Expected no byte to be read, returned %d (buffer content: %q)", n, buf)
		}
		if err == nil || err != io.EOF {
			t.Fatal("Expected EOF, got", err)
		}

		if err := f.Close(); err != nil {
			t.Fatal("Error closing file:", err)
		}
	})

	t.Run("create same file twice", func(t *testing.T) {
		mfs := NewMemFS()

		if err := mfs.CreateFile(filePath, []byte(fileData)); err != nil {
			t.Fatal("Error creating file:", err)
		}

		err := mfs.CreateFile(filePath, []byte("new content"))
		if err == nil {
			t.Fatal("Expected error creating file that already exists")
		}
		if expectErr := fs.ErrExist; !errors.Is(err, expectErr) {
			t.Fatalf("Expected error %q, got %q", expectErr, err)
		}
	})

	t.Run("open nonexistent file", func(t *testing.T) {
		mfs := NewMemFS()

		_, err := mfs.Open(filePath)
		if err == nil {
			t.Fatal("Expected error opening nonexistent file", err)
		}
		if expectErr := fs.ErrNotExist; !errors.Is(err, expectErr) {
			t.Fatalf("Expected error %q, got %q", expectErr, err)
		}
	})
}
