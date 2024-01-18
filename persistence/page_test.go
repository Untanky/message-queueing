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
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	WriteTo(io.Writer) (int64, error)
	ReadFrom(io.Reader) (int64, error)
}

type ssTablePageTest[Value ssTablePage] interface {
	newPage() Value
	setupPage(Value)
	compare(t *testing.T, a, b Value)
}

func TestSSTablePages(t *testing.T) {
	type testCase[Value ssTablePage] struct {
		name         string
		testLogic    ssTablePageTest[Value]
		expectedHash string
	}

	cases := []testCase[*tableHeader]{
		{name: "tableHeader", testLogic: tableHeaderTest{}, expectedHash: "pWENWZQonO0jeVSWKrl6Xhqen1S/psnOYGsbthbyA9w="},
	}

	random = rand.New(rand.NewSource(10))
	now = func() time.Time {
		return time.Date(2024, 10, 1, 14, 40, 0, 0, time.Local)
	}
	uuid.SetRand(random)

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			runTestCases(t, c.testLogic, c.expectedHash)
		})
	}
}

func runTestCases[Value ssTablePage](t *testing.T, test ssTablePageTest[Value], expectedHash string) {
	if test == nil {
		return
	}

	t.Run("Marshalling and Unmarshalling", func(tt *testing.T) {
		entity := test.newPage()
		test.setupPage(entity)

		bytes, err := entity.Marshal()
		if err != nil {
			tt.Errorf("err: expected nil; got %v", err)
		}

		actual := test.newPage()
		err = actual.Unmarshal(bytes)
		if err != nil {
			tt.Errorf("err: expected nil; got %v", err)
		}

		test.compare(tt, entity, actual)
	})

	t.Run("Writing and Reading", func(tt *testing.T) {
		handler := &testutils.SliceReadWriteSeeker{}

		entity := test.newPage()
		test.setupPage(entity)
		entity.WriteTo(handler)

		hash := sha256.New()
		hash.Write(handler.Data)
		hashBytes := hash.Sum(nil)

		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

		if hashBase64 != expectedHash {
			tt.Errorf("hashBytes: expected %v; got %v", expectedHash, hashBase64)
		}

		handler.Seek(0, io.SeekStart)
		actual := test.newPage()
		actual.ReadFrom(handler)

		test.compare(tt, entity, actual)
	})
}
