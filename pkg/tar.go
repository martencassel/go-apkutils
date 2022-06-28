package apkutils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

func WriteTarGz(filepath string, src io.ReadCloser, dst io.Writer) error {
	zr := gzip.NewWriter(dst)
	tw := tar.NewWriter(zr)
	fileInfo, err := os.Lstat(filepath)
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
	if err != nil {
		return err
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	data, err := os.Open(filepath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(tw, data); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}
