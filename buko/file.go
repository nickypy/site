package server

import (
	"io/fs"
	"log"
	"os"
	"path"
	fp "path/filepath"
	"sync"
)

func ensureDir(d string) {
	dir := path.Dir(d)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
}

func readFile(path string) []byte {
	contents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

func writeFile(filepath string, contents []byte) {
	filepath, err := fp.Abs(filepath)
	if err != nil {
		panic(err)
	}
	ensureDir(filepath)

	log.Default().Printf("Writing out %s", filepath)
	err = os.WriteFile(filepath, contents, fs.FileMode(0644))
	if err != nil {
		log.Default().Println(err)
	}
}

func CopyDirectory(source, dest string) error {
	var wg sync.WaitGroup

	err := fp.WalkDir(source, func(filepath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		dest := path.Join(dest, path.Base(filepath))

		wg.Go(func() {
			CopyFile(filepath, dest)
		})

		return nil
	})

	wg.Wait()
	return err
}

func CopyFile(source, dest string) {
	contents := readFile(source)
	writeFile(dest, contents)
}
