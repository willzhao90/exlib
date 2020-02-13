package mongo

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	pb "gitlab.com/sdce/protogo"
)

var (
	uuidPointer = reflect.TypeOf((*pb.UUID)(nil))
	uuidType    = reflect.TypeOf(pb.UUID{})
)

func GetRegistry() *bsoncodec.Registry {
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bson.PrimitiveCodecs{}.RegisterPrimitiveCodecs(rb)
	registerUUIDCodecs(rb)
	return rb.Build()
}

func registerUUIDCodecs(rb *bsoncodec.RegistryBuilder) {
	if rb == nil {
		panic(errors.New("argument to RegisterUUIDCodecs must not be nil"))
	}

	rb.
		RegisterEncoder(uuidPointer, bsoncodec.ValueEncoderFunc(uuidPointerEncoder)).
		RegisterDecoder(uuidPointer, bsoncodec.ValueDecoderFunc(uuidPointerDecoder)).
		RegisterEncoder(uuidType, bsoncodec.ValueEncoderFunc(uuidEncoder))
}

func uuidPointerDecoder(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {

	if !val.CanSet() || val.Kind() != uuidPointer.Kind() {
		return bsoncodec.ValueDecoderError{Name: "StringDecodeValue", Kinds: []reflect.Kind{uuidPointer.Kind()}, Received: val}
	}

	var uuid = val.Interface().(*pb.UUID)

	if vr.Type() == bsontype.Null {
		err := vr.ReadNull()
		val.Set(reflect.Zero(val.Type()))
		return err
	}

	if vr.Type() != bsontype.ObjectID {
		return fmt.Errorf("cannot decode %v into a UUID type", vr.Type())
	}
	oid, err := vr.ReadObjectID()

	if err != nil {
		return err
	}
	uuid.Bytes = make([]byte, len(oid))
	copy(uuid.Bytes, oid[:])
	return nil
}

// UUIDEncoder encoder
func uuidPointerEncoder(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if val.Kind() != uuidPointer.Kind() {
		return bsoncodec.ValueEncoderError{Name: "UUIDEncodeValue", Kinds: []reflect.Kind{uuidPointer.Kind()}, Received: val}
	}
	var out primitive.ObjectID

	var uuid = val.Interface().(*pb.UUID)

	if uuid == nil {
		vw.WriteNull()
		return nil
	}

	if len(uuid.Bytes) > 12 {
		return bsoncodec.ValueEncoderError{Name: "UUIDEncodeValue", Kinds: []reflect.Kind{uuidPointer.Kind()}, Received: val}
	}

	copy(out[:], uuid.Bytes[:12])
	return vw.WriteObjectID(out)
}

// UUIDEncoder encoder
func uuidEncoder(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if val.Kind() != uuidType.Kind() {
		return bsoncodec.ValueEncoderError{Name: "UUIDEncodeValue", Kinds: []reflect.Kind{uuidType.Kind()}, Received: val}
	}
	var out primitive.ObjectID

	var uuid = val.Interface().(pb.UUID)

	if len(uuid.Bytes) > 12 {
		return bsoncodec.ValueEncoderError{Name: "UUIDEncodeValue", Kinds: []reflect.Kind{uuidType.Kind()}, Received: val}
	}

	copy(out[:], uuid.Bytes[:12])
	return vw.WriteObjectID(out)
}
