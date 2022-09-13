package adapters

import (
	"context"
	"filippo.io/age"
	"filippo.io/age/armor"
	"fmt"
	"go-secretshelper/pkg/core"
	"io"
	"log"
	"strings"
)

// AgeEncryptTransformationType is the type identifier to be used in configuration files
const AgeEncryptTransformationType = "age"

// AgeEncryptTransformation is an transformation port adapter capable of
// encrypting a secret using the age encryption scheme.
type AgeEncryptTransformation struct {
	log *log.Logger
}

// AgeEncryptionTransformationSpec is the specification part for the configuration
type AgeEncryptionTransformationSpec struct {
	// Recipient is the public key of the age recipient. Must be a age1 X25519 type of recipient.
	Recipient string `yaml:"recipient" validate:"required"`
}

// NewAgeEncryptionTransformationSpec creates an AgeEncryptionTransformationSpec from a generic map
func NewAgeEncryptionTransformationSpec(in map[interface{}]interface{}) (AgeEncryptionTransformationSpec, error) {
	recipient, ok := in["recipient"]
	if !ok {
		return AgeEncryptionTransformationSpec{}, fmt.Errorf("recipient element is required in spec of age transform")
	}
	recipientStr, ok := recipient.(string)
	if !ok {
		return AgeEncryptionTransformationSpec{}, fmt.Errorf("recipient element must be a string in spec of age transform")
	}

	return AgeEncryptionTransformationSpec{
		Recipient: recipientStr,
	}, nil
}

// NewAgeEncryptTransformation returns a new instance of AgeEncrypt transformation
func NewAgeEncryptTransformation(log *log.Logger) *AgeEncryptTransformation {
	return &AgeEncryptTransformation{log: log}
}

// ProcessSecret takes the incoming secret, encrypts it for the recipient as per spec of the transformation and
// returns it as a new secret
func (aet *AgeEncryptTransformation) ProcessSecret(ctx context.Context, defaults *core.Defaults, secrets *core.Secrets, transformation *core.Transformation) (*core.Secret, error) {

	spec, err := NewAgeEncryptionTransformationSpec(transformation.Spec)
	if err != nil {
		return nil, err
	}

	recp, err := age.ParseX25519Recipient(spec.Recipient)
	if err != nil {
		return nil, err
	}

	// look up all input secrets from transformation, append
	// to one large string. make an input reader from this
	rawInput := ""
	for _, inVar := range *secrets {
		for _, inName := range transformation.Input {
			if inVar.Name == inName {
				rawInput = fmt.Sprintf("%s%s", rawInput, inVar.RawContent)
			}
		}
	}
	in := strings.NewReader(rawInput)

	// make an armored output writer as a string builder
	b := new(strings.Builder)
	a := armor.NewWriter(b)
	defer func() {
		if err := a.Close(); err != nil {
			aet.log.Printf("AgeTransformation: error closing writer: %v", err)
		}
	}()

	w, err := age.Encrypt(a, recp)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(w, in); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	if err := a.Close(); err != nil {
		return nil, err
	}

	return &core.Secret{
		Name:           transformation.Output,
		Type:           "transformed-by-age",
		RawContent:     []byte(b.String()),
		RawContentType: "application/octet-stream", //spec.ContentType,
	}, nil
}
