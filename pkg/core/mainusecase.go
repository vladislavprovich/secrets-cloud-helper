package core

import (
	"context"
	"errors"
	"fmt"
	"log"
)

// MainUseCaseImpl implements the UseCase interface
type MainUseCaseImpl struct {
	log *log.Logger
}

// NewMainUseCaseImpl creates a new main use case
func NewMainUseCaseImpl(l *log.Logger) UseCase {
	return &MainUseCaseImpl{
		log: l,
	}
}

// RetrieveSecret pulls a single secret from a vault and puts it into a Repository
func (m *MainUseCaseImpl) RetrieveSecret(ctx context.Context, factory Factory, defaults *Defaults, repository Repository, vault *Vault, secret *Secret) error {

	va := factory.NewVaultAccessor(vault.Type)
	if va == nil {
		return errors.New("internal error: unable to handle vault of given type")
	}

	updatedSecret, err := va.RetrieveSecret(ctx, defaults, vault, secret)
	if err != nil {
		return err
	}

	repository.Put(secret.Name, updatedSecret)

	return nil
}

// Transform applies transformation steps to repository
func (m *MainUseCaseImpl) Transform(ctx context.Context, factory Factory,
	defaults *Defaults, repository Repository, secrets *Secrets, transformation *Transformation) error {

	secretByName := make(map[string]*Secret)
	for _, secret := range *secrets {
		s, err := repository.Get(secret.Name)
		if err != nil {
			return err
		}
		secretByName[secret.Name] = s.(*Secret)
	}

	tr := factory.NewTransformation(transformation.Type)

	// collect all secrets that are required by the transformation
	in := make(Secrets, 0)
	for _, inputVarName := range transformation.Input {
		s, ex := secretByName[inputVarName]
		if !ex {
			return fmt.Errorf("transformation: input variable %s not found", inputVarName)
		}
		in = append(in, s)
	}

	m.log.Printf("Calling ProcessSecret %#v, %#v, %#v, %#v", ctx, defaults, in, transformation)
	updatedSecret, err := tr.ProcessSecret(ctx, defaults, &in, transformation)
	if err != nil {
		return err
	}

	// add secrets to config
	//*secrets = append(*secrets, updatedSecret)

	repository.Put(updatedSecret.Name, updatedSecret)

	return nil
}

// WriteToSink writes output a single sink by pulling it from the repository
func (m *MainUseCaseImpl) WriteToSink(ctx context.Context, factory Factory, defaults *Defaults, repository Repository, sink *Sink) error {

	// get secret to be written from repository
	repositoryContent, err := repository.Get(sink.Var)
	if err != nil {
		return err
	}

	var secret *Secret = repositoryContent.(*Secret)

	// get sik writer for type
	sw := factory.NewSinkWriter(sink.Type)
	if sw == nil {
		return errors.New("internal error: unable to handle sink of given type")
	}

	// write to sink and be done
	err = sw.Write(ctx, defaults, secret, sink)
	if err != nil {
		return err
	}

	return nil
}

// need at least one secret, from one vault going to one sink. If either is missing, we cannot proceed.
func (m *MainUseCaseImpl) dataMissing(vaults *Vaults, secrets *Secrets, sinks *Sinks) bool {
	return (secrets == nil || len(*secrets) == 0) || (vaults == nil || len(*vaults) == 0) || (sinks == nil || len(*sinks) == 0)
}

// Process runs the main use case
func (m *MainUseCaseImpl) Process(ctx context.Context, factory Factory, defaults *Defaults,
	vaults *Vaults, secrets *Secrets, transformations *Transformations, sinks *Sinks) error {

	if m.dataMissing(vaults, secrets, sinks) {
		return nil
	}

	repo := factory.NewRepository()

	m.log.Printf("Pulling secrets from vaults")
	for _, secret := range *secrets {
		vault := vaults.GetVaultByName(secret.VaultName)
		if vault == nil {
			return fmt.Errorf("no such vault: %s", secret.VaultName)
		}
		if err := m.RetrieveSecret(ctx, factory, defaults, repo, vault, secret); err != nil {
			return err
		}
	}

	// Applying transformations
	if transformations != nil && len(*transformations) > 0 {
		m.log.Printf("Applying transformations")
		for _, transformation := range *transformations {
			if err := m.Transform(ctx, factory, defaults, repo, secrets, transformation); err != nil {
				return err
			}
		}
	}

	// writing to all sinks
	m.log.Printf("Writing secrets to sinks")
	if sinks != nil {
		for _, sink := range *sinks {
			if err := m.WriteToSink(ctx, factory, defaults, repo, sink); err != nil {
				return err
			}
		}
	}

	return nil
}
