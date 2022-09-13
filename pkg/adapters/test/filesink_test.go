package test

import (
	"context"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"log"
	"reflect"
	"testing"
)

func TestFileSinkSpec(t *testing.T) {
	fss, err := adapters.NewFileSinkSpec(core.SinkSpec{
		"path":  "tmp",
		"mode":  440,
		"user":  0,
		"group": "-1",
	})
	if err != nil {
		t.Errorf("unexpected: %s", err)
	}
	if fss.Path != "tmp" {
		t.Errorf("expected tmp, got: %#v", fss.Path)

	}
	if fss.Mode == nil {
		t.Error("unexpected: got nil")
	} else {
		if *fss.Mode != 440 {
			t.Errorf("expected 440, got: %#v", fss)
		}
	}
	if fss.UserID == nil {
		t.Error("unexpected: got nil")
	} else {
		if *fss.UserID != 0 {
			t.Errorf("expected 0, got: %#v", *fss.UserID)
		}
	}
	if fss.GroupID == nil {
		t.Error("unexpected: got nil")
	} else {
		if *fss.GroupID != -1 {
			t.Errorf("expected -1, got: %#v", fss)
		}
	}

	fss, err = adapters.NewFileSinkSpec(core.SinkSpec{
		"path": "",
		"mode": "440",
	})
	if err != nil {
		t.Errorf("unexpected: %s", err)
	}
	if fss.Mode == nil {
		t.Error("unexpected: got nil")
	} else {
		if *fss.Mode != 440 {
			t.Errorf("expected 440, got: %#v", fss)
		}
	}
	if fss.UserID != nil {
		t.Errorf("expected nil, got %#v", fss.UserID)
	}
	if fss.GroupID != nil {
		t.Errorf("expected nil, got %#v", fss.UserID)
	}

	fss, err = adapters.NewFileSinkSpec(core.SinkSpec{
		"mode": true,
	})
	if err == nil {
		t.Error("Expected error, got none.")
	}
}

func TestFileSink(t *testing.T) {
	secrets := &core.Secrets{
		&core.Secret{
			Name:       "test",
			Type:       "secret",
			VaultName:  "test",
			RawContent: []byte("s3cr3t"),
		},
	}
	p := "test.dat"
	sinks := &core.Sinks{
		&core.Sink{
			Type: "mock",
			Var:  "test",
			Spec: core.SinkSpec{
				"path": p,
				"mode": "440",
			},
		},
	}

	fs := afero.NewMemMapFs()

	sink := adapters.NewFileSink(log.Default(), fs)
	err := sink.Write(context.TODO(), &core.Defaults{}, (*secrets)[0], (*sinks)[0])
	if err != nil {
		t.Errorf("Unexpected: %s", err)
	}

	//
	fi, err := fs.Stat(p)
	if err != nil {
		t.Errorf("Unexpected: %s", err)
	}
	if fi.Size() != int64(len((*secrets)[0].RawContent)) {
		t.Errorf("Invalid size")
	}

	if fi.Mode().Perm() != 440 {
		t.Errorf("Expected mode 440, got: %#v", fi.Mode().Perm())
	}

	raw, err := afero.ReadFile(fs, p)
	if err != nil {
		t.Errorf("Unexpected: %s", err)
	}

	cmp := reflect.DeepEqual(raw, (*secrets)[0].RawContent)
	if !cmp {
		t.Errorf("Invalid content")
	}
}
