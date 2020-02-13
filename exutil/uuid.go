package exutil

import (
	"encoding/hex"
	"fmt"

	pb "gitlab.com/sdce/protogo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewUUID generates random protobuf uuid
func NewUUID() *pb.UUID {
	oid := primitive.NewObjectID()
	return &pb.UUID{Bytes: oid[:]}
}

// UUIDtoA encodes uuid to string
func UUIDtoA(uuid *pb.UUID) string {
	if uuid == nil {
		return ""
	}
	return hex.EncodeToString(uuid.Bytes)
}

// AtoUUID decodes string to UUID
func AtoUUID(txt string) (*pb.UUID, error) {
	uuid, err := hex.DecodeString(txt)
	if err != nil {
		return nil, fmt.Errorf("cannot decode uuid %v: %v", txt, err)
	}
	if len(uuid) != 12 {
		return nil, fmt.Errorf("expecting 12 bytes UUID, got %v", uuid)
	}
	return &pb.UUID{Bytes: uuid}, nil
}

// ObjectIDtoUUID converts mongo ObjectID to UUID
func ObjectIDtoUUID(oid primitive.ObjectID) *pb.UUID {
	return &pb.UUID{Bytes: oid[:]}
}

// UUIDtoObjectID converts UUID to mongo ObjectID
func UUIDtoObjectID(uuid *pb.UUID) (oid primitive.ObjectID, err error) {
	if len(uuid.Bytes) != 12 {
		err = fmt.Errorf("length of uuid exceeds 12 bytes")
		return
	}

	if n := copy(oid[:], uuid.Bytes[:]); n != 12 {
		err = fmt.Errorf("Expecting to write 12 bytes, only written %v bytes", n)
	}
	return
}

// UUID is alias of pb.UUID
type UUID pb.UUID

func (id *UUID) String() string {
	if id == nil {
		return ""
	}
	return UUIDtoA((*pb.UUID)(id))
}
