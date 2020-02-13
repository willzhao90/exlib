package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorToRpcError(err error) error {
	switch err {
	case mongo.ErrNoDocuments:
		return status.Error(codes.NotFound, "Docs not found.")
	case mongo.ErrClientDisconnected,
		mongo.ErrEmptySlice,
		mongo.ErrInvalidIndexValue,
		mongo.ErrMissingResumeToken,
		mongo.ErrMultipleIndexDrop,
		mongo.ErrNilCursor,
		mongo.ErrNilDocument,
		mongo.ErrNonStringIndexName,
		mongo.ErrUnacknowledgedWrite,
		mongo.ErrWrongClient:
		s := status.New(codes.Internal, "database unavailable")
		return s.Err()
	}
	return err
}
