package api

import (
	"github.com/meiti-x/go-transactional-msg/pb"
	"google.golang.org/grpc"
	"log"
	"os"
	"path/filepath"
)

type fsGRPCApi struct {
	pb.UnimplementedFSServer
	filePath string
}

func NewFsGRPCApi(filepath string) pb.FSServer {
	return &fsGRPCApi{filePath: filepath}
}

func (s *fsGRPCApi) Upload(req grpc.ClientStreamingServer[pb.Chunk, pb.UploadResponse]) error {
	var (
		filename string
		content  []byte
	)

	for {
		chunk, err := req.Recv()

		if err != nil {
			log.Println(err.Error())

			req.SendAndClose(&pb.UploadResponse{
				Status: pb.UploadStatus_FAILED,
			})
		}

		if len(filename) == 0 {
			filename = chunk.Filename
		}

		content = append(content, chunk.Data)

		if chunk.Done {
			break
		}
	}

	if err := os.WriteFile(filepath.Join(s.filePath, filename), content, os.ModePerm); err != nil {
		log.Println(err.Error())

		req.SendAndClose(&pb.UploadResponse{
			Status: pb.UploadStatus_FAILED,
		})
	}

	req.SendAndClose(&pb.UploadResponse{
		Status: pb.UploadStatus_SUCCESS,
	})
	return nil
}
