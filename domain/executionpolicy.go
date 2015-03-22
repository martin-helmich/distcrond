package domain

import "errors"

const (
	POLICY_ALL = "all"
	POLICY_ANY = "any"
)

type ExecutionPolicyJson struct {
	Hosts string
	HostList []string
	Roles []string
}

type ExecutionPolicy struct {
	Hosts string
	HostList []string
	Roles []string
}

func NewExecutionPolicyFromJson(json ExecutionPolicyJson) (ExecutionPolicy, error) {
	return ExecutionPolicy {
		Hosts: json.Hosts,
		HostList: json.HostList,
		Roles: json.Roles,
	}, nil
}

func (p ExecutionPolicy) IsValid() error {
	if p.Hosts == POLICY_ALL || p.Hosts == POLICY_ANY {
		if len(p.HostList) == 0 && len(p.Roles) == 0 {
			return errors.New("Either 'HostList' or 'Roles' must have at least one entry")
		}
	} else {
		return errors.New("'Hosts' must be 'all' or 'any'")
	}

	return nil
}
