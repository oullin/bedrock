package filesystem

import (
	"bufio"
	"encoding/json"
	"iter"
	"os"
)

// Get reads the entire contents of a file.
func (f *Filesystem) Get(path string) ([]byte, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return data, nil
}

// JSON reads a file and unmarshals its JSON contents into v.
func (f *Filesystem) JSON(path string, v any) error {
	data, err := f.Get(path)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// SharedGet reads a file's contents while holding a shared (read) lock.
func (f *Filesystem) SharedGet(path string) ([]byte, error) {
	file, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	defer file.Close()

	if err := lockShared(file); err != nil {
		return nil, err
	}

	defer func() {
		_ = unlockFile(file)
	}()

	return os.ReadFile(path)
}

// Lines returns an iterator that yields each line of the file.
func (f *Filesystem) Lines(path string) (iter.Seq[string], error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return func(yield func(string) bool) {
		file, err := os.Open(path)

		if err != nil {
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}, nil
}
