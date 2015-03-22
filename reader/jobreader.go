package reader

import (
	"log"
	"os"
	"errors"
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"github.com/martin-helmich/distcrond/domain"
)

type JobReader struct {
	validationConfig domain.JobValidationConfig
	receiver JobReceiver
}

type JobReceiver interface {
	AddJob(domain.Job)
}

func NewJobReader(validationConfig domain.JobValidationConfig, receiver JobReceiver) *JobReader {
	reader := new(JobReader)
	reader.validationConfig = validationConfig
	reader.receiver = receiver

	return reader
}

func (r JobReader) ReadFromDirectory(directory string) error {
	log.Println("Reading configuration")

	var walk filepath.WalkFunc = func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}

		if file.Name()[0] == '.' {
			log.Println("Skipping", path)
			return nil
		}

		fileContents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		wrapError := func(err error) error {
			return errors.New(fmt.Sprintf("Error parsing file %s: %s", path, err))
		}

		name := strings.Replace(file.Name(), ".json", "", 1)
		jobJson := domain.JobJson{}

		err = json.Unmarshal(fileContents, &jobJson)
		if err != nil {
			return wrapError(err)
		}

		job, mappingErr := domain.NewJobFromJson(name, jobJson)
		if mappingErr != nil {
			return wrapError(mappingErr)
		}

		log.Printf("Read job from %s: %s\n", file.Name(), job)

		if validErr := job.IsValid(r.validationConfig); validErr != nil {
			return wrapError(validErr)
		}

		r.receiver.AddJob(job)

		return nil
	}

	if err := filepath.Walk(directory, walk); err != nil {
		return err
	}

	return nil
}
