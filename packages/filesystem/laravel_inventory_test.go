package filesystem_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/filesystem"
)

// FilesystemTest::testGetRetrievesFiles
// FilesystemTest::testPutStoresFiles
// FilesystemTest::testLines
// FilesystemTest::testLinesThrowsExceptionNonexisitingFile
// FilesystemTest::testReplaceCreatesFile
// FilesystemTest::testReplaceInFileCorrectlyReplaces
// FilesystemTest::testReplaceWhenUnixSymlinkExists
// FilesystemTest::testSetChmod
// FilesystemTest::testGetChmod
// FilesystemTest::testDeleteRemovesFiles
// FilesystemTest::testPrependExistingFiles
// FilesystemTest::testPrependNewFiles
// FilesystemTest::testMissingFile
// FilesystemTest::testDeleteDirectory
// FilesystemTest::testDeleteDirectoryReturnFalseWhenNotADirectory
// FilesystemTest::testCleanDirectory
// FilesystemTest::testFilesMethod
// FilesystemTest::testCopyDirectoryReturnsFalseIfSourceIsntDirectory
// FilesystemTest::testCopyDirectoryMovesEntireDirectory
// FilesystemTest::testMoveDirectoryMovesEntireDirectory
// FilesystemTest::testMoveDirectoryMovesEntireDirectoryAndOverwrites
// FilesystemTest::testMoveDirectoryReturnsFalseWhileOverwritingAndUnableToDeleteDestinationDirectory
// FilesystemTest::testGetThrowsExceptionNonexisitingFile
// FilesystemTest::testJsonReturnsDecodedJsonData
// FilesystemTest::testJsonReturnsNullIfJsonDataIsInvalid
// FilesystemTest::testAppendAddsDataToFile
// FilesystemTest::testMoveMovesFiles
// FilesystemTest::testNameReturnsName
// FilesystemTest::testExtensionReturnsExtension
// FilesystemTest::testBasenameReturnsBasename
// FilesystemTest::testDirnameReturnsDirectory
// FilesystemTest::testTypeIdentifiesFile
// FilesystemTest::testTypeIdentifiesDirectory
// FilesystemTest::testSizeOutputsSize
// FilesystemTest::testMimeTypeOutputsMimeType
// FilesystemTest::testIsWritable
// FilesystemTest::testIsReadable
// FilesystemTest::testIsDirEmpty
// FilesystemTest::testGlobFindsFiles
// FilesystemTest::testAllFilesFindsFiles
// FilesystemTest::testDirectoriesFindsDirectories
// FilesystemTest::testAllDirectoriesFindsDirectories
// FilesystemTest::testMakeDirectory
// FilesystemTest::testSharedGet
// FilesystemTest::testCopyCopiesFileProperly
// FilesystemTest::testHasSameHashChecksFileHashes
// FilesystemTest::testIsFileChecksFilesProperly
// FilesystemTest::testFilesMethodReturnsFileInfoObjects
// FilesystemTest::testAllFilesReturnsFileInfoObjects
// FilesystemTest::testHashWithDefaultValue
// FilesystemTest::testHash
// FilesystemTest::testLastModifiedReturnsTimestamp
// FilesystemTest::testFileCreationAndContentVerification
// FilesystemTest::testDirectoryOperationsWithSubdirectories
func TestUpstreamFilesystemInventoryEquivalents(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	file := filepath.Join(dir, "docs", "example.txt")

	if err := fs.Put(file, []byte("World"), 0o640); err != nil {
		t.Fatal(err)
	}

	if got, err := fs.Get(file); err != nil || string(got) != "World" {
		t.Fatalf("expected file contents, got %q err=%v", string(got), err)
	}

	if info, err := os.Stat(file); err != nil || info.Mode().Perm() != 0o640 {
		t.Fatalf("expected chmod mode 0640, got info=%v err=%v", info, err)
	}

	if err := fs.Prepend(file, []byte("Hello ")); err != nil {
		t.Fatal(err)
	}

	if err := fs.Append(file, []byte("!")); err != nil {
		t.Fatal(err)
	}

	if got, err := fs.SharedGet(file); err != nil || string(got) != "Hello World!" {
		t.Fatalf("expected shared contents, got %q err=%v", string(got), err)
	}

	if err := fs.ReplaceInFile("World", "Bedrock", file); err != nil {
		t.Fatal(err)
	}

	if got, err := fs.Get(file); err != nil || string(got) != "Hello Bedrock!" {
		t.Fatalf("expected replaced contents, got %q err=%v", string(got), err)
	}

	linesFile := filepath.Join(dir, "lines.txt")

	if err := fs.Put(linesFile, []byte("one\ntwo\nthree")); err != nil {
		t.Fatal(err)
	}

	lines, err := fs.Lines(linesFile)

	if err != nil {
		t.Fatal(err)
	}

	var collected []string

	for line := range lines {
		collected = append(collected, line)
	}

	if len(collected) != 3 || collected[0] != "one" || collected[2] != "three" {
		t.Fatalf("unexpected lines: %#v", collected)
	}

	if _, err := fs.Lines(filepath.Join(dir, "missing.txt")); !errors.Is(err, filesystem.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for missing lines file, got %v", err)
	}

	jsonFile := filepath.Join(dir, "payload.json")

	if err := fs.Put(jsonFile, []byte(`{"name":"Taylor"}`)); err != nil {
		t.Fatal(err)
	}

	var payload map[string]string

	if err := fs.JSON(jsonFile, &payload); err != nil || payload["name"] != "Taylor" {
		t.Fatalf("unexpected JSON payload: %#v err=%v", payload, err)
	}

	badJSON := filepath.Join(dir, "bad.json")

	if err := fs.Put(badJSON, []byte("{bad json")); err != nil {
		t.Fatal(err)
	}

	if err := fs.JSON(badJSON, &payload); err == nil {
		t.Fatal("expected invalid JSON to return an error")
	}

	if !fs.Exists(file) || fs.Missing(file) || !fs.IsFile(file) {
		t.Fatal("expected file existence helpers to identify the file")
	}

	if fs.Name(file) != "example" || fs.Extension(file) != "txt" || fs.Basename(file) != "example.txt" || fs.Dirname(file) != filepath.Dir(file) {
		t.Fatal("unexpected path metadata")
	}

	if typ, err := fs.Type(filepath.Dir(file)); err != nil || typ != "dir" {
		t.Fatalf("expected dir type, got %q err=%v", typ, err)
	}

	if typ, err := fs.Type(file); err != nil || typ != "file" {
		t.Fatalf("expected file type, got %q err=%v", typ, err)
	}

	if size, err := fs.Size(file); err != nil || size == 0 {
		t.Fatalf("expected non-zero file size, got %d err=%v", size, err)
	}

	if _, err := fs.LastModified(file); err != nil {
		t.Fatal(err)
	}

	if mime, err := fs.MimeType(file); err != nil || mime == "" {
		t.Fatalf("expected MIME type, got %q err=%v", mime, err)
	}

	if !fs.IsReadable(file) || !fs.IsWritable(file) {
		t.Fatal("expected file to be readable and writable")
	}

	hashDefault, err := fs.Hash(file)

	if err != nil {
		t.Fatal(err)
	}

	hashSHA1, err := fs.Hash(file, "sha1")

	if err != nil {
		t.Fatal(err)
	}

	if hashDefault == "" || hashSHA1 == "" || hashDefault == hashSHA1 {
		t.Fatalf("unexpected hashes: md5=%q sha1=%q", hashDefault, hashSHA1)
	}

	copyPath := filepath.Join(dir, "copy.txt")

	if err := fs.Copy(file, copyPath); err != nil {
		t.Fatal(err)
	}

	same, err := fs.HasSameHash(file, copyPath)

	if err != nil || !same {
		t.Fatalf("expected copied file to share hash, same=%v err=%v", same, err)
	}

	movePath := filepath.Join(dir, "moved.txt")

	if err := fs.Move(copyPath, movePath); err != nil {
		t.Fatal(err)
	}

	if fs.Exists(copyPath) || !fs.Exists(movePath) {
		t.Fatal("expected file move to remove source and create target")
	}

	if err := fs.Delete(movePath); err != nil {
		t.Fatal(err)
	}

	if fs.Exists(movePath) {
		t.Fatal("expected delete to remove file")
	}

	if _, err := fs.Get(filepath.Join(dir, "missing.txt")); !errors.Is(err, filesystem.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for missing get, got %v", err)
	}

	nested := filepath.Join(dir, "nested")

	if err := fs.MakeDirectory(filepath.Join(nested, "a", "b")); err != nil {
		t.Fatal(err)
	}

	if empty, err := fs.IsEmptyDirectory(filepath.Join(nested, "a", "b"), false); err != nil || !empty {
		t.Fatalf("expected empty directory, empty=%v err=%v", empty, err)
	}

	if err := fs.Put(filepath.Join(nested, "a", "b", "file.txt"), []byte("content")); err != nil {
		t.Fatal(err)
	}

	if err := fs.Put(filepath.Join(nested, "root.txt"), []byte("content")); err != nil {
		t.Fatal(err)
	}

	if files, err := fs.Files(nested); err != nil || len(files) != 1 {
		t.Fatalf("expected one direct file, got %#v err=%v", files, err)
	}

	if files, err := fs.AllFiles(nested); err != nil || len(files) != 2 {
		t.Fatalf("expected two recursive files, got %#v err=%v", files, err)
	}

	if dirs, err := fs.Directories(nested); err != nil || len(dirs) != 1 {
		t.Fatalf("expected one direct directory, got %#v err=%v", dirs, err)
	}

	if dirs, err := fs.AllDirectories(nested); err != nil || len(dirs) != 2 {
		t.Fatalf("expected two recursive directories, got %#v err=%v", dirs, err)
	}

	if matches, err := fs.Glob(filepath.Join(nested, "*.txt")); err != nil || len(matches) != 1 {
		t.Fatalf("expected one glob match, got %#v err=%v", matches, err)
	}

	copiedDir := filepath.Join(dir, "copied")

	if err := fs.CopyDirectory(nested, copiedDir); err != nil {
		t.Fatal(err)
	}

	movedDir := filepath.Join(dir, "moved-dir")

	if err := fs.MoveDirectory(copiedDir, movedDir, true); err != nil {
		t.Fatal(err)
	}

	if err := fs.CleanDirectory(movedDir); err != nil {
		t.Fatal(err)
	}

	if empty, err := fs.IsEmptyDirectory(movedDir, false); err != nil || !empty {
		t.Fatalf("expected cleaned directory to be empty, empty=%v err=%v", empty, err)
	}

	if err := fs.DeleteDirectory(movedDir); err != nil {
		t.Fatal(err)
	}

	if err := fs.CopyDirectory(file, filepath.Join(dir, "not-dir")); !errors.Is(err, filesystem.ErrNotDirectory) {
		t.Fatalf("expected ErrNotDirectory for copy directory from file, got %v", err)
	}

	if err := fs.DeleteDirectory(file); !errors.Is(err, filesystem.ErrNotDirectory) {
		t.Fatalf("expected ErrNotDirectory for delete directory from file, got %v", err)
	}
}
