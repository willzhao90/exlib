package exutil

import (
	"bytes"
	"encoding/hex"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestUUIDGenLen(t *testing.T) {
	generated := NewUUID()
	log.Printf("Generated: %x", generated.Bytes)
	if len(generated.Bytes) != 12 {
		log.Println("Generated uuid should have 12 bytes for mongo.")
		t.Fail()
	}
}

func TestUUIDtoA(t *testing.T) {
	generated := NewUUID()
	log.Printf(UUIDtoA(generated))
	uuid, err := AtoUUID(UUIDtoA(generated))
	if err != nil {

		t.Fail()
		return
	}
	log.Printf("Converted %x", uuid.Bytes)
	if !bytes.Equal(uuid.Bytes, generated.Bytes) {
		t.Fail()
	}
}

func TestUUIDtoOID(t *testing.T) {
	generated := NewUUID()
	oid, err := UUIDtoObjectID(generated)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	uuid := ObjectIDtoUUID(oid)
	if !bytes.Equal(uuid.Bytes, generated.Bytes) {
		t.Fail()
	}
	log.Printf("From %x, through %v, got %x back.", generated.Bytes, oid, uuid.Bytes)
}

func TestUUID_String(t *testing.T) {
	src := "5cf860e16391ed9def4d4fb9"
	bts, _ := hex.DecodeString(src)
	id := &UUID{bts}
	if id.String() != src {
		t.Errorf("expect %v , got %v", src, id)
	}
}
