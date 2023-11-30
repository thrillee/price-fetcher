package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"

	"github.com/thrillee/price-fetcher/proto"
	"google.golang.org/grpc"
)

func GRPCServerlistenAndServe(listenAddr string, svc PriceFetcher) error {
	fmt.Printf("GRPC Server (%s) started ...", listenAddr)
	grpcPriceFetcher := NewGRPCPriceFetcherServer(svc)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	proto.RegisterPriceFetcherServer(server, grpcPriceFetcher)

	return server.Serve(ln)
}

type GRPCPriceFetcherServer struct {
	svc PriceFetcher
	proto.UnimplementedPriceFetcherServer
}

func NewGRPCPriceFetcherServer(svc PriceFetcher) *GRPCPriceFetcherServer {
	return &GRPCPriceFetcherServer{
		svc: svc,
	}
}

func (s *GRPCPriceFetcherServer) FetchPrice(ctx context.Context, req *proto.PriceRequest) (*proto.PriceResponse, error) {
	reqid := rand.Intn(10000)
	ctx = context.WithValue(ctx, "requestID", reqid)

	price, err := s.svc.FetchPrice(ctx, req.Ticker)
	if err != nil {
		return nil, err
	}

	resp := &proto.PriceResponse{
		Ticker: req.Ticker,
		Price:  float32(price),
	}

	return resp, nil
}
