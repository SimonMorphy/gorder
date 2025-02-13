package client

import (
	"context"
	"errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net"
	"time"

	"github.com/SimonMorphy/gorder/common/discovery"
	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
	"github.com/SimonMorphy/gorder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	if !waitForStockGrpcClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("order grpc is not available")
	}
	serviceName := viper.GetString("stock.service-name")
	serviceAddr, err := discovery.GetServiceAddr(ctx, serviceName)
	if err != nil {
		return nil, func() error {
			return nil
		}, nil
	}
	addr := serviceAddr
	if addr == "" {
		logrus.Warn("empty grpc addr for stock grpc")
	}
	options, err := grpcDialOptions(addr)
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	conn, err := grpc.NewClient(addr, options...)
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func waitForOrderGrpcClient(timeout time.Duration) bool {
	logrus.Infof("waiting for order grpc clinet,duration is %v", timeout)
	return waitFor(viper.GetString("order.grpc-addr"), timeout)
}

func waitForStockGrpcClient(timeout time.Duration) bool {
	logrus.Infof("waiting for stock grpc clinet,duration is %v", timeout)
	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
}

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	if !waitForOrderGrpcClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("order grpc is not available")
	}
	serviceName := viper.GetString("order.service-name")
	serviceAddr, err := discovery.GetServiceAddr(ctx, serviceName)
	if err != nil {
		return nil, func() error {
			return nil
		}, nil
	}
	addr := serviceAddr
	if addr == "" {
		logrus.Warn("empty grpc addr for order grpc")
	}
	options, err := grpcDialOptions(addr)
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	conn, err := grpc.NewClient(addr, options...)
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func grpcDialOptions(_ string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}, nil
}

func waitFor(addr string, timeout time.Duration) bool {
	ch := make(chan struct{})
	after := time.After(timeout)

	go func() {
		for {
			select {
			case <-after:
				return
			default:
			}
			_, err := net.Dial("tcp", addr)
			if err == nil {
				close(ch)
				return
			}
			time.Sleep(2 * time.Second)
		}
	}()

	select {
	case <-ch:
		return true
	case <-after:
		return false
	}
}
