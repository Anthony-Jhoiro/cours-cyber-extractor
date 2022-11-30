package commons

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
)

func BuildFile(id uint32) (string, error) {
	tmpFileDir := fmt.Sprintf("results/%d", id)
	fileInfos, err := os.ReadDir(tmpFileDir)
	if err != nil {
		return "", err
	}

	file, err := os.Create(fmt.Sprintf("results/%d.raw", id))

	defer file.Close()

	sort.Slice(fileInfos, func(i, j int) bool {
		f1, err := strconv.Atoi(fileInfos[i].Name())
		if err != nil {
			return false
		}

		f2, err := strconv.Atoi(fileInfos[j].Name())
		if err != nil {
			return false
		}
		return f1 < f2
	})

	for _, fileInfo := range fileInfos {
		if !fileInfo.Type().IsDir() {
			content, err := os.ReadFile(path.Join(tmpFileDir, fileInfo.Name()))
			if err != nil {
				return "", err
			}
			_, err = file.Write(content)
			if err != nil {
				return "", err
			}
		}
	}

	os.RemoveAll(tmpFileDir)

	return file.Name(), nil
}

func WriteByteFile(id uint32, seq uint32, payload []byte) {
	log.Printf("Recieve [%d] - No %d - %d bytes\n", id, seq, len(payload))

	os.Mkdir(fmt.Sprintf("results/%d", id), 0777)
	file, err := os.Create(fmt.Sprintf("results/%d/%d", id, seq))
	if err != nil {
		panic(err)
	}
	file.Write(payload)

	file.Close()
}
