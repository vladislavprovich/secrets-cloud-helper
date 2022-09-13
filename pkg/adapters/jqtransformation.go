package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itchyny/gojq"
	"go-secretshelper/pkg/core"
	"log"
	"strings"
)

// JQTransformationType is the type string to be used within configuration files for the jq transformation
const JQTransformationType = "jq"

// JQTransformation is a transformation adapter that uses jq to transform the input
type JQTransformation struct {
	log *log.Logger
}

// NewJQTransformation creates a new JQTransformation
func NewJQTransformation(log *log.Logger) *JQTransformation {
	return &JQTransformation{log: log}
}

// JQTransformationSpec contains the specification items for the jq transformation
type JQTransformationSpec struct {
	// Query is the jq query string
	Query string `json:"query" validate:"required"`

	// If raw is set, the result is rendered as a string and not as a json. This works
	// only for simple values, not for arrays or objects.
	Raw bool `json:"raw" validate:"required"`

	// Content Type of rendered output (default: application/json for raw=false, text/plain for raw=true)
	ContentType string

	query *gojq.Query
}

// NewJQTransformationSpec creates a new JQTransformationSpec
func NewJQTransformationSpec(in map[interface{}]interface{}) (JQTransformationSpec, error) {
	var spec JQTransformationSpec

	qs, ok := in["q"]
	if !ok {
		return JQTransformationSpec{}, fmt.Errorf("q (query) element is required")
	}
	q, ok := qs.(string)
	if !ok {
		return JQTransformationSpec{}, fmt.Errorf("q (query) element must be a string")
	}
	var err error
	spec.query, err = gojq.Parse(q)
	if err != nil {
		return JQTransformationSpec{}, err
	}

	raw, ok := in["raw"]
	if ok {
		rawBool, ok := raw.(bool)
		if ok {
			spec.Raw = rawBool
		}
	}

	spec.Query = q
	var contentType = "application/json"
	cnt, ok := in["contentType"]
	if ok {
		cntStr, ok := cnt.(string)
		if ok {
			contentType = cntStr
		}
	}
	spec.ContentType = contentType

	return spec, nil
}

// ProcessSecret returns a new secret as the result of a json query process
func (t *JQTransformation) ProcessSecret(ctx context.Context,
	defaults *core.Defaults, in *core.Secrets, transformation *core.Transformation) (*core.Secret, error) {

	spec, err := NewJQTransformationSpec(transformation.Spec)
	if err != nil {
		return nil, err
	}

	bIn := new(strings.Builder)
	bOut := new(strings.Builder)
	for _, inVar := range *in {
		bIn.Write(inVar.RawContent)
	}

	var input interface{}
	if err := json.Unmarshal([]byte(bIn.String()), &input); err != nil {
		return nil, err
	}

	var contentType = ""

	iter := spec.query.RunWithContext(ctx, input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}

		if spec.Raw {
			contentType = "text/plain"
			bStr, ok := v.(string)
			if ok {
				bOut.WriteString(bStr)
			} else {
				t.log.Printf("Warning: Unable to convert result of jq query to string: %v", v)
			}
		} else {
			contentType = "application/json"
			// interpret jq result as-is, as json, write to buffer
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			bOut.Write(b)
		}
	}

	// check for content type override
	if spec.ContentType != "" {
		contentType = spec.ContentType
	}

	res := &core.Secret{
		Name:           transformation.Output,
		Type:           "transformed-by:jq",
		RawContent:     []byte(bOut.String()),
		RawContentType: contentType,
	}

	return res, nil
}
