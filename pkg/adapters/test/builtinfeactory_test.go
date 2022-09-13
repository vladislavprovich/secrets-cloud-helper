package test

import (
	"github.com/spf13/afero"
	"go-secretshelper/pkg/adapters"
	"io/ioutil"
	"log"
	"testing"
)

func TestBuiltinFactory(t *testing.T) {
	fs := afero.NewMemMapFs()
	bif := adapters.NewBuiltinFactory(log.New(ioutil.Discard, "", 0), fs)

	m := make(map[string]struct{})
	for _, st := range bif.SinkTypes() {
		m[st] = struct{}{}
	}

	if _, ok := m[adapters.FileSinkType]; !ok {
		t.Errorf("%s is missing", adapters.FileSinkType)
	}

	sw := bif.NewSinkWriter(adapters.FileSinkType)
	if sw == nil {
		t.Error("Unexcepted: nil")
	}
}
