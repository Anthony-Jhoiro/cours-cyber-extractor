package commons

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
)

//BuildFile build files located under results/{id} and concatenate them to produce a single binary file
func BuildFile(id uint32) (string, error) {
	// List all temporary files
	tmpFileDir := fmt.Sprintf("results/%d", id)
	fileInfos, err := os.ReadDir(tmpFileDir)
	if err != nil {
		return "", err
	}

	// Create the destination file
	file, err := os.Create(fmt.Sprintf("results/%d.raw", id))

	// Close the result file at the end of the function
	defer file.Close()

	// Sort the temporary files by name
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

	// For each temporary file, add its content to the reult file
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

	// Delete the temporary file directory
	os.RemoveAll(tmpFileDir)

	return file.Name(), nil
}

//WriteByteFile Write bytes into a file which will be located under results/{id} and be name after the seq parameter
func WriteByteFile(id uint32, seq uint32, payload []byte) {
	log.Printf("Recieve [%d] - No %d - %d bytes\n", id, seq, len(payload))

	// Create the base directories
	os.Mkdir("results", 0777)
	os.Mkdir(fmt.Sprintf("results/%d", id), 0777)

	// Create the file
	file, err := os.Create(fmt.Sprintf("results/%d/%d", id, seq))
	if err != nil {
		panic(err)
	}

	// Add the content into the file
	file.Write(payload)

	// Close the file
	file.Close()
}

//HandleStopFile build the file from all the temporary packets
func HandleStopFile(id uint32) {
	log.Printf("Recieved Stop file\n")

	// Build the files
	fileName, err := BuildFile(id)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}

	log.Printf("Recieved File %s", fileName)
}
