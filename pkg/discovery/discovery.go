package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	//register the service
	Register(ctx context.Context, servicename string, hostport string, instanceId string) error
	//deregister Service
	Deregister(ctx context.Context, servicename string, instanceId string) error

	//Recieve Service Addresss

	ServiceAddresses(ctx context.Context, servicename string) ([]string, error)

	//Peroidic call to report the status
	ReportHealthyState(instanceId string, serviceName string) error
}

// ErrNotFound is returned when no service addresses are
// found.
var ErrNotFound = errors.New("no service addresses found")

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.
		NewSource(time.Now().UnixNano())).Int())
}
