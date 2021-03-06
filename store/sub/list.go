package sub

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	sep = string(filepath.Separator)
)

// mkStoreWalkerFunc create a func to walk a (sub)store, i.e. list it's content
func mkStoreWalkerFunc(alias, folder string, fn func(...string)) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && path != folder {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		if path == folder {
			return nil
		}
		if path == filepath.Join(folder, GPGID) {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		s := strings.TrimPrefix(path, folder+sep)
		s = strings.TrimSuffix(s, ".gpg")
		if alias != "" {
			s = alias + sep + s
		}
		fn(s)
		return nil
	}
}

// List will list all entries in this store
func (s *Store) List(prefix string) ([]string, error) {
	lst := make([]string, 0, 10)
	addFunc := func(in ...string) {
		lst = append(lst, in...)
	}

	err := filepath.Walk(s.path, mkStoreWalkerFunc(prefix, s.path, addFunc))
	return lst, err
}
