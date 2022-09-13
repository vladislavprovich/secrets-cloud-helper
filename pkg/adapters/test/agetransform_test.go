package test

import (
	"context"
	"filippo.io/age"
	"filippo.io/age/armor"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestAgeEncryptTransformationSpec(t *testing.T) {
	in := core.TransformationSpec{
		"recipient": "age1njkx5t9tcc4gq7c53zzy4sfjq0fscm5uzt5vek5pj2khehcpsfsqwzq9jy",
	}

	spec, err := adapters.NewAgeEncryptionTransformationSpec(in)
	if err != nil {
		t.Errorf("Error creating spec: %s", err)
	}
	if spec.Recipient != "age1njkx5t9tcc4gq7c53zzy4sfjq0fscm5uzt5vek5pj2khehcpsfsqwzq9jy" {
		t.Errorf("Expected recipient to be present, got nil")
	}

	spec, err = adapters.NewAgeEncryptionTransformationSpec(core.TransformationSpec{})
	if err == nil {
		t.Errorf("Expected error creating spec")
	}
}

func TestAgeEncryptTransformation(t *testing.T) {
	secrets := &core.Secrets{
		{
			Name:       "s1",
			RawContent: []byte("s3cr3t"),
		},
	}
	transformation := &core.Transformation{
		Input:  []string{"s1"},
		Output: "s1-enc",
		Type:   "age",
		Spec: core.TransformationSpec{
			"recipient": "age1njkx5t9tcc4gq7c53zzy4sfjq0fscm5uzt5vek5pj2khehcpsfsqwzq9jy",
		},
	}

	tt := adapters.NewAgeEncryptTransformation(log.New(ioutil.Discard, "", 0))
	s, err := tt.ProcessSecret(context.TODO(), &core.Defaults{}, secrets, transformation)
	if err != nil {
		t.Errorf("Error processing template transformation: %s", err)
	}

	if s == nil {
		t.Errorf("Expected secret to be present, got nil")
	}

	// age-decrypt using identity
	const ageIdentity = "AGE-SECRET-KEY-15AZV7RD7N8V8KJKJ8WMUNV49JP8V5MG4FF8RRC8W4Z689E7FLSAQ4DCZHJ"
	r := strings.NewReader(ageIdentity)
	identities, err := age.ParseIdentities(r)
	if err != nil {
		t.Errorf("Error parsing identities: %s", err)
	}
	in := armor.NewReader(strings.NewReader(string(s.RawContent)))
	decryptReader, err := age.Decrypt(in, identities...)
	if err != nil {
		t.Errorf("Error decrypting: %s", err)
	}
	sb := new(strings.Builder)
	if _, err := io.Copy(sb, decryptReader); err != nil {
		t.Errorf("Error decrypting: %s", err)
	}

	if reflect.DeepEqual(sb.String(), string((*secrets)[0].RawContent)) == false {
		t.Errorf("Expected secret to be %s, got %s", sb.String(), string((*secrets)[0].RawContent))
	}

}
