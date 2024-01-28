package persistence

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"message-queueing/testutils"
	"testing"
	"time"
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

	random = rand.New(rand.NewSource(10))
	now = func() time.Time {
		return time.Date(2024, 10, 1, 14, 40, 0, 0, time.Local)
	}
	uuid.SetRand(random)

	t.Run("tableHeader", func(tt *testing.T) {
		runTestCases[*tableHeader](tt, tableHeaderTest{}, "ld21tWrGxJBxD4927I0V+m1MU04VltLwpTB3z2U2mtk=")
	})
	t.Run("dataPage", func(tt *testing.T) {
		runTestCases[*dataPage](tt, dataPageTest{}, "gkpvioStVRtGudjcGih/jQBmQ0CC3KJbqijTSgxUeek=")
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
