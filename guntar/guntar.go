package guntar

import (
	"archive/tar"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func ArchiveDirectory(dirName string, outputPath string) error {
	if fileExists(outputPath) {
		// do not overwrite
		log.Fatalf("Archive output file %s already exists", outputPath)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Error creating output file %s, error: %v", outputPath, err)
	}
	defer f.Close()

	// the thing that actually writes the tar file
	writer := tar.NewWriter(f)
	defer writer.Close()

	// traverse *all* the files in the directory
	// we do not use WalkDir because we actually the FileInfo for each file
	// TODO: return error instead of dying with log.Fatalf
	filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Could not walk file (%s): %v", path, err)
		}

		canonPath := path
		if strings.HasPrefix(path, dirName) {
			// remove directory path prefix
			// for example, /tmp/something/dir/file becomes /dir/file
			// this is to convert host's path to container/image path
			canonPath = path[len(dirName):]
		} else {
			log.Fatalf("Paths of files found in directory (%s) somehow do not have the dir as prefix: %s", dirName, path)
		}

		if info.IsDir() {
			// directories must have a / at the end
			canonPath += "/"
		}

		// if the file is a symlink, this will be the linked path
		linkPath := canonPath
		isSymLink := info.Mode() & os.ModeSymlink == os.ModeSymlink
		if isSymLink {
			linkPath, err = os.Readlink(path)
			if err != nil {
				log.Fatalf("Could not get link path (%s): %v", path, err)
			}
		}

		// log.Println(canonPath, info.Name(), info.IsDir())

		// make tar file header
		header, err := tar.FileInfoHeader(info, linkPath)
		if err != nil {
			log.Fatalf("Could not make file header (%s): %v", path, err)
		}
		// the header's Name will be path in the host (/tmp/something/dir/file)
		// we need the container's path (/dir/file)
		header.Name = canonPath

		// write header
		if err := writer.WriteHeader(header); err != nil {
			log.Fatalln("Error writing tar header: ", err)
		}

		// if it's a regular file, we write the contents
		if !info.IsDir() && !isSymLink {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("Could not read file (%s): %v", path, err)
			}
			if _, err := writer.Write(data); err != nil {
				log.Fatalf("Error writing file/directory %s into tar file: %v", canonPath, err)
			}
		}

		return nil
	})

	// ensure everything is written
	writer.Flush()

	// done!
	if err := writer.Close(); err != nil {
		log.Fatalf("Error closing tar file %s: %v", outputPath, err)
	}

	return nil
}
