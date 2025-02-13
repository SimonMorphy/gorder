package server

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net"

	glogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	gtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	glogrus.ReplaceGrpcLogger(logrus.NewEntry(logrus.StandardLogger()))
}

func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
	addr := viper.Sub(serviceName).GetString("grpc-addr")
	if addr == "" {
		//TODO: Warning Log
		addr = viper.GetString("fallback-grpc-addr")
	}
	logrus.Info("grpc server run on ", addr)
	RunGPRCServerOnAddr(addr, registerServer)
}

func RunGPRCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpcSrv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			// Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
			gtags.UnaryServerInterceptor(gtags.WithFieldExtractor(gtags.CodeGenRequestFieldExtractor)),
			glogrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc.ChainStreamInterceptor(
			gtags.StreamServerInterceptor(gtags.WithFieldExtractor(gtags.CodeGenRequestFieldExtractor)),
			glogrus.StreamServerInterceptor(logrusEntry),
		),
	)
	registerServer(grpcSrv)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Panic(err)
	}
	if err = grpcSrv.Serve(listen); err != nil {
		logrus.Panic(err)
	}
}
