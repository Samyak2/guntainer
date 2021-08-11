package guntar

import (
	"archive/tar"
	"errors"
	"log"
	"os"
	"path/filepath"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func ArchiveDirectory(dirName string, outputPath string) error {
	if fileExists(outputPath) {
		log.Fatalf("Archive output file %s already exists", outputPath)
	}
	f, err := os.Create(outputPath)
	defer f.Close()
	if err != nil {
		log.Fatalf("Error creating output file %s, error: %v", outputPath, err)
	}
	writer := tar.NewWriter(f)
	defer writer.Close()

	filepath.WalkDir(dirName, func(path string, info os.DirEntry, err error) error {
		// TODO: take each file/directory and tar it

		log.Println(path, info.Name(), info.IsDir())
		header := &tar.Header{
			Name: path,
			Linkname: "",  // TODO

			Size: 0, // TODO
			Mode: 0, // TODO
		}

		if err := writer.WriteHeader(header); err != nil {
			log.Fatalln("Error writing tar header: ", err)
		}

		// TODO: read file here and pass bytes
		if _, err := writer.Write([]byte{}); err != nil {
			// TODO: replace path here with (path with prefix removed)
			log.Fatalf("Error writing file/directory %s into tar file: %v", path, err)
		}

		return nil
	})

	writer.Flush()

	if err := writer.Close(); err != nil {
		log.Fatalf("Error closing tar file %s: %v", outputPath, err)
	}

	return nil
}
