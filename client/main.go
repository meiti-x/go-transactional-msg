package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/meiti-x/go-transactional-msg/pb"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var filePath = flag.String("file", "", "file to send")

func main() {
	flag.Parse()

	if *filePath == "" {
		log.Fatal("No file provided to upload!")
	}

	client, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	fsClient := pb.NewFSClient(client)

	data, err := os.ReadFile(*filePath)

	if err != nil {
		log.Fatal(err)
	}

	fileName := filepath.Base(*filePath)
	stream, err := fsClient.Upload(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	start, end := 0, 500
	chunkSize := 500
	for {
		if end >= len(data) {
			end = len(data)
		}

		done := end == len(data)

		dataToSend := data[start:end]

		stream.Send(&pb.Chunk{
			Data:     dataToSend,
			FileName: fileName,
			Done:     done,
		})

		if done {
			break
		}

		start, end = end, end+chunkSize

	}

	resp, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Upload status: ", resp.Status)
}
