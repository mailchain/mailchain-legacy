package datastoretest

import "github.com/spf13/afero"

// ReadDir returns the files in the given directory
func ReadDir(fs afero.Fs, root string) []string {
	files := make([]string, 0, 1)

	f, err := fs.Open(root)
	if err != nil {
		panic(err)
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		panic(err)
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return files
}
