// Package fs provides filesystem-related functions.
package fs

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/etherlabsio/errors"
	"github.com/etherlabsio/pkg/env"
)

const (
	// DirMode is the default permission used when creating directories
	DirMode = 0755
	// FileMode is the default permission used when creating files
	FileMode = 0644
)

// Gopath will return the current GOPATH as set by environment variables and
// will fall back to ~/go if a GOPATH is not set.
func Gopath() string {
	home := env.String("HOME", "~/")
	return env.String("GOPATH", filepath.Join(home, "go"))
}

// CopyDir is a utility to assist with copying a directory from src to dest.
// Note that directory permissions are not maintained, but the permissions of
// the files in those directories are.
func CopyDir(src, dest string) error {
	const op = errors.Op("CopyDir")
	dir, err := os.Open(src)
	if err != nil {
		return errors.WithOp(err, op)
	}
	if err := os.MkdirAll(dest, DirMode); err != nil {
		return errors.WithOp(err, op)
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return errors.WithOp(err, op)
	}
	for _, file := range files {
		srcptr := filepath.Join(src, file.Name())
		dstptr := filepath.Join(dest, file.Name())
		if file.IsDir() {
			if err := CopyDir(srcptr, dstptr); err != nil {
				return errors.WithOp(err, op)
			}
		} else {
			if err := CopyFile(srcptr, dstptr); err != nil {
				return errors.WithOp(err, op)
			}
		}
	}
	return nil
}

// CopyFile is a utility to assist with copying a file from src to dest.
// Note that file permissions are maintained.
func CopyFile(src, dest string) error {
	const op = errors.Op("CopyFile")
	source, err := os.Open(src)
	if err != nil {
		return errors.WithOp(err, op)
	}
	defer source.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return errors.WithOp(err, op)
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, source)
	if err != nil {
		return errors.WithOp(err, op)
	}
	sourceinfo, err := os.Stat(src)
	if err != nil {
		return errors.WithOp(err, op)
	}

	return os.Chmod(dest, sourceinfo.Mode())
}

// UntarBundle will untar a source tar.gz archive to the supplied destination
func UntarBundle(destination string, source string) error {
	const op = errors.Op("UntarBundle")
	f, err := os.Open(source)
	if err != nil {
		return errors.New("open download source", err, op, errors.IO)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return errors.New("create gzip reader from "+source, errors.IO, err, op)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("reading tar file", err, op, errors.IO)
		}

		path := filepath.Join(filepath.Dir(destination), header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return errors.New("creating directory for tar file: "+path, errors.IO, op, err)
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return errors.New("open file "+path, err, op, errors.IO)
		}
		defer file.Close()
		if _, err := io.Copy(file, tr); err != nil {
			msg := fmt.Sprintf("copy tar %s to destination %s", header.FileInfo().Name(), path)
			return errors.New(msg, err, op, errors.IO)
		}
	}
	return nil
}

// AbsFilePath resolves ENV vars specified in string for file paths
func AbsFilePath(inPath string) string {
	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}
	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}
	if p, err := filepath.Abs(inPath); err == nil {
		return filepath.Clean(p)
	}
	return ""
}
