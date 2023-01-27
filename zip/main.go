package main

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

func main() {
	body, err := os.ReadFile("./data.zip")
	if err != nil {
		panic(err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(body))
}

func prettyStr(source []byte, size uint) string {
	pretty := "[ "
	for i, b := range source {

		if math.Mod(float64(i), float64(size)) == 0 {
			pretty += "\n"
		}

		pretty += fmt.Sprintf("%d, ", b)
	}

	return pretty[:len(pretty)-2] + "\n]"
}

func write(source, target string) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func unzip(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	return io.ReadAll(f)
}
