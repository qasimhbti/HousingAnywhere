package main

import (
	"context"
	"io"
	"log"

	"github.com/grpcengine/pb"
	"google.golang.org/grpc"
)

const port = ":9000"

func main() {
	/*creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}*/
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("localhost"+port, opts...)
	if err != nil {
		log.Fatalf("failed to dial %v", err)
	}
	defer conn.Close()
	client := pb.NewHousingAnywhereSeviceClient(conn)

	/*res, err := client.CalcBreakEven(context.Background(), &pb.Breakeven{
		Downpayment:           5000,
		Intratemortgage:       7.5,
		Propertytaxes:         30,
		Propertytransfertaxes: 15,
		Yearstolive:           15,
		Totalpropertycost:     12000,
		Monthlyrent:           600,
	})
	if err != nil {
		log.Fatalf("rpc call %v:", err)
	}
	log.Println("BreakEven :", res.IsOk)*/

	/*stream, err := client.CalcBreakEven(context.Background(), &pb.BreakEvenRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.Employee)
	}*/

	CalculateBreakEven(client)
}

func CalculateBreakEven(client pb.HousingAnywhereSeviceClient) {
	properties := []pb.Property{
		pb.Property{
			PropertyID:            1,
			Downpayment:           5000,
			Intratemortgage:       7.5,
			Propertytaxes:         30.67,
			Propertytransfertaxes: 15.77,
			Yearstolive:           15,
			Totalpropertycost:     12000,
			Monthlyrent:           600,
			IsBreakEven:           false,
		},
		pb.Property{
			PropertyID:            2,
			Downpayment:           5000,
			Intratemortgage:       3.5,
			Propertytaxes:         5.56,
			Propertytransfertaxes: 8.89,
			Yearstolive:           15,
			Totalpropertycost:     12000,
			Monthlyrent:           1600,
			IsBreakEven:           false,
		},
	}

	stream, err := client.CalcBreakEven(context.Background())
	if err != nil {
		log.Fatalf("rpc call %v:", err)
	}
	doneCh := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				doneCh <- struct{}{}
				break
			}
			if err != nil {
				log.Fatalf("stream receive err %v", err)
			}
			log.Printf("PropertyID :%d, BreakEven : %t", res.Property.PropertyID, res.Property.IsBreakEven)
		}
	}()
	for _, property := range properties {
		err := stream.Send(&pb.BreakEvenRequest{
			Property: &property,
		})
		if err != nil {
			log.Fatalf("stream send error %v:", err)
		}
	}
	_ = stream.CloseSend()
	<-doneCh
}
