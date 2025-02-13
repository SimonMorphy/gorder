package discovery

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

type Registry interface {
	Registry(ctx context.Context, instanceID, serviceName, hostPort string) error
	Deregister(ctx context.Context, instanceID, serviceName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceID, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	x := rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Int()
	return fmt.Sprintf("%s-%d", serviceName, x)
}
