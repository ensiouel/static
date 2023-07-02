package static

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ensiouel/basket-contract/gen/go/static/v1"
	"github.com/ensiouel/static/internal/domain/static"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
)

type Server struct {
	staticService static.Service
	maxFileSize   int

	pb_static.UnimplementedStaticServer
}

func NewStaticServer(staticService static.Service, maxFileSize int) *Server {
	return &Server{staticService: staticService, maxFileSize: maxFileSize}
}

func (server *Server) Upload(stream pb_static.Static_UploadServer) error {
	temp, err := os.CreateTemp("temp", "file")
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("error creating temp file: %v", err))
	}

	defer func() {
		temp.Close()

		err = os.Remove(temp.Name())
		if err != nil && !os.IsNotExist(err) {
			log.Printf("error removing temp file: %v", err)
			return
		}
	}()

	fileSize := 0
	hash := sha256.New()

	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			return status.Error(codes.Internal, fmt.Sprintf("error receiving file: %v", err))
		}

		data := recv.GetData()

		_, err = temp.Write(data)
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("error writing file: %v", err))
		}

		hash.Write(data)

		fileSize += len(data)
		if fileSize > server.maxFileSize {
			return status.Error(codes.InvalidArgument, "file too large")
		}
	}

	if fileSize == 0 {
		return status.Error(codes.InvalidArgument, "file is empty")
	}

	err = temp.Close()
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("error closing file: %v", err))
	}

	hashsum := hex.EncodeToString(hash.Sum(nil))

	filename, err := server.staticService.Upload(stream.Context(), temp, hashsum)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("error uploading file: %v", err))
	}

	return stream.SendAndClose(&pb_static.UploadResponse{
		SourceId: filename,
	})
}

func (server *Server) Download(request *pb_static.DownloadRequest, stream pb_static.Static_DownloadServer) error {
	sourceID := request.GetSourceId()

	if sourceID == "" {
		return status.Error(codes.InvalidArgument, "source id is required")
	}

	file, err := server.staticService.Download(stream.Context(), sourceID)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("error downloading file: %v", err))
	}

	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			return status.Error(codes.Internal, fmt.Sprintf("error reading bytes: %v", err))
		}

		err = stream.Send(&pb_static.DownloadResponse{Data: buf[:n]})
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("error sending bytes: %v", err))
		}
	}

	return nil
}

func (server *Server) Delete(ctx context.Context, request *pb_static.DeleteRequest) (*pb_static.DeleteResponse, error) {
	sourceID := request.GetSourceId()

	if sourceID == "" {
		return nil, status.Error(codes.InvalidArgument, "source id is required")
	}

	err := server.staticService.Delete(ctx, sourceID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error deleting file: %v", err))
	}

	return &pb_static.DeleteResponse{}, nil
}
