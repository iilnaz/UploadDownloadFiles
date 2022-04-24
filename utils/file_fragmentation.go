package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	fileChunk int64 = 1 * (1 << 20)
)

var realSizeOfChunks map[string]int64

func createFileDir(dirName string) {
	err := os.MkdirAll(fmt.Sprintf("storage/%s", dirName), 0777)
	if err != nil {
		panic(err)
	}
}

func ChunkFile(file multipart.File, header *multipart.FileHeader) error {
	realSizeOfChunks = make(map[string]int64)

	createFileDir(header.Filename)

	fileName := strings.Split(header.Filename, ".")

	defer file.Close()

	fileSize := header.Size

	partNum := math.Ceil(float64(fileSize) / float64(fileChunk))

	for i := 0; i < int(partNum); i++ {
		partBuf := make([]byte, fileChunk)

		nameOfChunk := fmt.Sprintf("storage/%s/%s_%d", header.Filename, fileName[0], i)
		countOfRead, err := file.Read(partBuf)
		realSizeOfChunks[nameOfChunk] = int64(countOfRead)

		tmp, err := os.Create(nameOfChunk)
		if err != nil {
			tmp.Close()
			err = errors.Wrap(err, "error while creating chunk")
			return err
		}
		err = ioutil.WriteFile(nameOfChunk, partBuf, os.ModeAppend)
		if err != nil {
			err = errors.Wrap(err, "error while writing to chunk")
		}
		tmp.Close()
	}
	return nil
}

func GlueFiles(w http.ResponseWriter, fileName string) error {
	files, _ := ioutil.ReadDir(fmt.Sprintf("storage/%s", fileName))
	countOfFiles := len(files)
	chunkName := strings.Split(fileName, ".")

	for i := 0; i < countOfFiles; i++ {
		nameOfFile := fmt.Sprintf("storage/%s/%s_%d", fileName, chunkName[0], i)
		file, err := os.Open(nameOfFile)
		if err != nil {
			err = errors.Wrap(err, "file not found")
		}
		if realSizeOfChunks[nameOfFile] != fileChunk {
			buf, _ := os.ReadFile(nameOfFile)
			buf = buf[0:realSizeOfChunks[nameOfFile]]
			tmp, err := ioutil.TempFile(fmt.Sprintf("storage/%s", fileName), "tmp")
			if err != nil {
				err = errors.Wrap(err, "error while creating new file")
				return err
			}
			if _, err = io.Copy(w, tmp); err != nil {
				err = errors.Wrap(err, "error while copy from tmp file")
				return err
			}
			if err = os.Remove(tmp.Name()); err != nil {
				err = errors.Wrap(err, "error while removing tmp file")
				return err
			}
		} else {
			if _, err = io.Copy(w, file); err != nil {
				err = errors.Wrap(err, "error while copy from chunk")
			}
		}
		file.Close()
	}
	return nil
}
