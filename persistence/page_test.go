package persistence

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"message-queueing/testutils"
	"testing"
)

type ssTablePage interface {
	WriteTo(io.Writer) (int64, error)
	ReadFrom(io.Reader) (int64, error)
}

type ssTablePageTest[Value ssTablePage] interface {
	newPage() Value
	setupPage(Value)
	compare(t *testing.T, a, b Value)
}

func TestSSTablePages(t *testing.T) {
	uuid.SetRand(rand.New(rand.NewSource(10)))

	t.Run("tableHeader", func(tt *testing.T) {
		runTestCases[*tableHeader](tt, tableHeaderTest{}, "vjX1qoot5KRAGjleMfK6X5FNo2hUrVPhYvxcowU6D88=")
	})
	t.Run("dataPage", func(tt *testing.T) {
		runTestCases[*dataPage](tt, dataPageTest{}, "aKbzyPmRicWQjUkWjDjEw91jJr9bAZiPSBZ61HDQays=")
	})
	t.Run("secondaryIndexPage", func(tt *testing.T) {
		runTestCases[*secondaryIndexPage[int]](tt, secondaryIndexPageTest{}, "aXDx/YYCDYvaVRBbOpMXNJfwtnQmky8mozVcczGwJ3o=")
	})
	t.Run("inMemorySecondaryIndex", func(tt *testing.T) {
		runTestCases[*inMemorySecondaryIndex[int]](tt, secondaryIndexTest{}, "VXyw4z/wDOE+DPr08fogHH5dZW4u4CJTfmRDsBAYh3A=")
	})
}

func runTestCases[Value ssTablePage](t *testing.T, test ssTablePageTest[Value], expectedHash string) {
	handler := &testutils.SliceReadWriteSeeker{}

	entity := test.newPage()
	test.setupPage(entity)
	entity.WriteTo(handler)

	hash := sha256.New()
	hash.Write(handler.Data)
	hashBytes := hash.Sum(nil)

	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	if hashBase64 != expectedHash {
		t.Errorf("hashBytes: expected %v; got %v", expectedHash, hashBase64)
	}

	handler.Seek(0, io.SeekStart)
	actual := test.newPage()
	actual.ReadFrom(handler)

	test.compare(t, entity, actual)
}
