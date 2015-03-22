package domain

import "errors"

type JobOwnerJson struct {
	Name string
	EmailAddress string
}

type JobOwner struct {
	Name string
	EmailAddress string
}

func NewJobOwnerFromJson(json JobOwnerJson) (JobOwner, error) {
	return JobOwner {
		Name: json.Name,
		EmailAddress: json.EmailAddress,
	}, nil
}

func (o JobOwner) IsValid() error {
	if len(o.Name) == 0 {
		return errors.New("Owner name must not be empty")
	}

	if len(o.EmailAddress) == 0 {
		return errors.New("Owner email address must not be empty")
	}

	return nil
}
