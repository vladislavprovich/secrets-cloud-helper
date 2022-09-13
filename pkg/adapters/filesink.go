package adapters

import (
	"context"
	"errors"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/core"
	"log"
	"os"
	"strconv"
)

// FileSinkType is the valid type name for a file sink
const FileSinkType = "file"

// FileSinkSpec is a specialisation of the SinkSpec interface for file sink
type FileSinkSpec struct {
	Path    string  `yaml:"path" validate:"required"`
	Mode    *uint32 `yaml:"mode,omitempty" validate:"required"`
	UserID  *int    `yaml:"user,omitempty" validate:"required"`
	GroupID *int    `yaml:"group,omitempty" validate:"required"`
}

// FileSink is a file-based sink endpoint, where secrets are written to files
type FileSink struct {
	log *log.Logger
	fs  afero.Fs
}

// NewFileSink creates a new FileSink, based on given Afero file system and a logger
func NewFileSink(log *log.Logger, fs afero.Fs) *FileSink {
	return &FileSink{
		log: log,
		fs:  fs,
	}
}

// NewFileSinkSpec creates a FileSinkSpec struct from abstract map
func NewFileSinkSpec(in map[interface{}]interface{}) (FileSinkSpec, error) {
	var res FileSinkSpec

	var defaultMode uint32 = 400
	res.Mode = &defaultMode

	v, ex := in["path"]
	if !ex {
		return res, errors.New("must provide a path element for a file sink spec")
	}
	res.Path = v.(string)

	v, ex = in["mode"]
	if ex {
		vn, err := stringOrIntToI(v)
		if err != nil {
			return res, errors.New("mode parameter in file sink spec must be string or integer")
		}
		var vn2 uint32 = uint32(vn)
		res.Mode = &vn2
	}
	v, ex = in["user"]
	if ex {
		vn, err := stringOrIntToI(v)
		if err != nil {
			return res, errors.New("user parameter in file sink spec must be string or integer")
		}
		res.UserID = &vn
	}
	v, ex = in["group"]
	if ex {
		vn, err := stringOrIntToI(v)
		if err != nil {
			return res, errors.New("group parameter in file sink spec must be string or integer")
		}
		res.GroupID = &vn
	}

	return res, nil
}

func stringOrIntToI(v interface{}) (int, error) {
	var vn int
	var err error

	vs, ok := v.(string)
	if ok {
		vn, err = strconv.Atoi(vs)
		if err != nil {
			return vn, err
		}
	} else {
		vn, ok = v.(int)
		if !ok {
			return vn, errors.New("neither string nor int")
		}
	}

	return vn, nil
}

// Write writes secret to sink, sets owner and mode if given by spec
func (s *FileSink) Write(ctx context.Context, defaults *core.Defaults, secret *core.Secret, sink *core.Sink) error {

	spec, err := NewFileSinkSpec(sink.Spec)
	if err != nil {
		return err
	}

	f, err := s.fs.OpenFile(spec.Path, os.O_WRONLY|os.O_CREATE, os.FileMode(*spec.Mode))
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := f.Write(secret.RawContent)
	if err != nil {
		return err
	}
	if n != len(secret.RawContent) {
		return errors.New("invalid number of bytes")
	}

	uid := -1
	gid := -1
	if spec.UserID != nil {
		uid = *spec.UserID
	}
	if spec.GroupID != nil {
		gid = *spec.GroupID
	}
	if err := s.fs.Chown(spec.Path, uid, gid); err != nil {
		if err2 := s.fs.Remove(spec.Path); err2 != nil {
			s.log.Printf("Unable to chown %s to %d/%d, AND unable to delete file afterwards.", spec.Path, uid, gid)
		}

		return err
	}

	s.log.Printf("Written secret \"%s\" to file %s, mode %d, chown-ed %d:%d\n", secret.Name, spec.Path, os.FileMode(*spec.Mode), uid, gid)

	return nil
}
