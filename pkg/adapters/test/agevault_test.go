package test

import (
	"context"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"log"
	"os"
	"reflect"
	"testing"
)

// Key for age container testing
//# created: 2021-11-04T18:20:18+01:00
//# public key: age1dfamnuh6cwvk7c4p3nrlr027tm0urk5qqh49tq2udmxhzkltgayst7kuf3
//AGE-SECRET-KEY-1KGV83XLPAU7HVC3TS7WE5GS2QG5ECYH97W9ZPPGWAYUHCTYNEGVS8NMFQP

// echo '{ "test": "s3cr3t" }' | age -e -r age1dfamnuh6cwvk7c4p3nrlr027tm0urk5qqh49tq2udmxhzkltgayst7kuf3 -a
const ageVaultFile = `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBRTEM2NE9CQVNMSlRLSlhH
TXZIS1lDbjBIbWhOaDVsQjBuRklBTFpWNG5FCmFvZEJOY20wZThGRHFVVC9rUTFE
aUJ6MEYyeXpIM0VRRjUzVUY1MnJnK3cKLS0tIGcvUUVXREhHYS9hNVlNNFcxWlVh
MTUweXhlNGQyVzgwdVpVZ1Z3U2V6YncK+NlZ8k8Zv/BuyhLfRYu0o1SJF6klFSLk
Ev/8vDZnPjkXXLuwpXGWukIlEwms/SwA9fw86p0=
-----END AGE ENCRYPTED FILE-----
`
const ageIdentity = "AGE-SECRET-KEY-1KGV83XLPAU7HVC3TS7WE5GS2QG5ECYH97W9ZPPGWAYUHCTYNEGVS8NMFQP"

func setupAgeFiles(fs afero.Fs) error {
	f, err := fs.OpenFile("vault.age", os.O_WRONLY|os.O_CREATE, 0400)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(ageVaultFile))
	if err != nil {
		return err
	}

	f2, err := fs.OpenFile("identity.age", os.O_WRONLY|os.O_CREATE, 0400)
	if err != nil {
		return err
	}
	defer f2.Close()

	_, err = f2.Write([]byte(ageIdentity))
	if err != nil {
		return err
	}

	return nil
}

func TestAgeVault(t *testing.T) {
	fs := afero.NewMemMapFs()

	// place age-encrypted file into fs
	err := setupAgeFiles(fs)
	if err != nil {
		t.Error(err)
		return
	}

	av := adapters.NewAgeVault(log.New(os.Stdout, "***", 0), fs)
	if av == nil {
		t.Error("unexpected: nil")
	}

	vaults := &core.Vaults{
		&core.Vault{
			Name: "test",
			Type: "age-file",
			Spec: core.VaultSpec{
				"path":     "vault.age",
				"identity": "identity.age",
			},
		},
	}
	secrets := &core.Secrets{
		&core.Secret{
			Name:      "test",
			Type:      "secret",
			VaultName: "test",
		},
		&core.Secret{
			Name:      "nosuchsecret",
			Type:      "secret",
			VaultName: "test",
		},
	}

	// retrieve the secret
	res, err := av.RetrieveSecret(context.TODO(), &core.Defaults{}, (*vaults)[0], (*secrets)[0])
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if res == nil {
		t.Error("Unexpected nil")
	} else {
		if !reflect.DeepEqual(res.RawContent, []byte("s3cr3t")) {
			t.Error("Unexpected RawContent")
		}
	}

	_, err = av.RetrieveSecret(context.TODO(), &core.Defaults{}, (*vaults)[0], (*secrets)[1])
	if err == nil {
		t.Error("Unexpected nil")
	}

}
