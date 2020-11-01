package main

import (
	"io"
	"log"
	"net"

	"github.com/grpcengine/pb"
	"google.golang.org/grpc"
)

const port = ":9000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	/*creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}*/
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	pb.RegisterHousingAnywhereSeviceServer(s, new(housingAnywhereServer))
	log.Println("Starting server on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

type housingAnywhereServer struct {
	pb.UnimplementedHousingAnywhereSeviceServer
}

func (s *housingAnywhereServer) CalcBreakEven(stream pb.HousingAnywhereSevice_CalcBreakEvenServer) error {
	for {
		prop, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		totalExpensesOnRent := prop.Property.Monthlyrent * (prop.Property.Yearstolive * 12)
		totalExpensesOnBuy := prop.Property.Downpayment + (float32(prop.Property.Totalpropertycost)-prop.Property.Downpayment)*prop.Property.Intratemortgage + float32(prop.Property.Totalpropertycost)*prop.Property.Propertytaxes + float32(prop.Property.Totalpropertycost)*prop.Property.Propertytransfertaxes
		log.Println("Rent :", totalExpensesOnRent)
		log.Println("Buy :", totalExpensesOnBuy)
		if float32(totalExpensesOnRent) >= totalExpensesOnBuy {
			prop.Property.IsBreakEven = true
		}
		log.Println(prop.Property)
		_ = stream.Send(&pb.BreakEvenResponse{Property: prop.Property})
	}
	return nil
}
