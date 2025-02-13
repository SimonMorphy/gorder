package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/SimonMorphy/gorder/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error {
			return nil
		}, err
	}
	id := GenerateInstanceID(serviceName)
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
	if err := registry.Registry(ctx, id, serviceName, grpcAddr); err != nil {
		logrus.Errorf("fail to register to consul, reason: %v", err)
		return func() error {
			return nil
		}, nil
	}
	go func() {
		for {
			if err := registry.HealthCheck(id, serviceName); err != nil {
				logrus.Panic("no heartbeat reply from consul, service:%s ,err : %v", serviceName, err)
			}
			time.Sleep(time.Second)
		}
	}()
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        grpcAddr,
	}).Info("registered to consul")
	return func() error {
		return registry.Deregister(ctx, id, serviceName)
	}, nil
}

func GetServiceAddr(ctx context.Context, name string) (string, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}
	discover, err := registry.Discover(ctx, name)
	if err != nil {
		return "", err
	}
	if len(discover) == 0 {
		return "", fmt.Errorf("got empty target service:%s from consul", name)
	}
	i := rand.Intn(len(discover))
	logrus.Infof("Discovered %d service of %s , addr=%v", len(discover), name, discover)
	return discover[i], nil
}
