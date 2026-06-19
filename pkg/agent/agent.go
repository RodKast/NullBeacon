package agent

import (
	"fmt"

	"github.com/google/uuid"
)

type Agent struct {
	ID       string
	Username string
	Hostname string
	Address  string
}

func (a *Agent) String() string {
	return fmt.Sprintf("%s:%s", a.Username, a.Hostname)
}

func NewAgent(username, hostname, address string) *Agent {
	return &Agent{
		ID:       uuid.New().String(),
		Username: username,
		Hostname: hostname,
		Address:  address,
	}
}
