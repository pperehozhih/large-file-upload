package upload

import (
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pperehozhih/large-file-upload"
)

type Server struct {
	storage pb.Manager
}

func NewServer(storage pb.Manager) Server {
	return Server{
		storage: storage,
	}
}

func (s Server) UploadFile(stream pb.UploadService_UploadServer) error {
	name := "some-unique-name.png"
	file := storage.NewFile(name)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := s.storage.Store(file); err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			return stream.SendAndClose(&uploadpb.UploadResponse{Name: name})
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := file.Write(req.GetChunk()); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
