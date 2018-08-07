package grpc

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/9d77v/go-lib/clients/etcd"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/naming"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	gn "google.golang.org/grpc/naming"
)

//RegisterService register service to etcd and keep alive
func RegisterService(etcdCli *etcd.Client, serviceName string, profile string, hostPort string, ttl int64) {
	r := &naming.GRPCResolver{Client: etcdCli.Client}
	ticker := time.NewTicker(time.Duration(ttl) * time.Second)
	for {
		lease, err := etcdCli.Grant(context.TODO(), ttl)
		if err != nil {
			log.Println("lease failed")
		}
		err = r.Update(context.TODO(), serviceName+"-"+profile, gn.Update{Op: gn.Add, Addr: hostPort}, clientv3.WithLease(lease.ID))
		log.Println("register node :", serviceName+"-"+profile, hostPort, err)
		select {
		case <-ticker.C:
		}
	}
}

//NewClientConn new grpc client connection use etcd balancer
func NewClientConn(etcdCli *etcd.Client, serviceName string, profile string) (*grpc.ClientConn, error) {
	r := &naming.GRPCResolver{Client: etcdCli.Client}
	b := grpc.RoundRobin(r)
	return grpc.Dial(serviceName+"-"+profile,
		grpc.WithBalancer(b),
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
		)))
}

//CloseClient close clients' connections
type CloseClient func()

//SignalHandler check signal for grpceful stop
func SignalHandler(server *grpc.Server, closeClient CloseClient) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <-ch
		log.Printf("signal: %v", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Println("stop")
			signal.Stop(ch)
			log.Println("close clients' connections")
			closeClient()
			server.GracefulStop()
			log.Printf("graceful shutdown")
			return
		}
	}
}
