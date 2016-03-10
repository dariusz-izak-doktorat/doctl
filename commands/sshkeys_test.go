package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bryanl/doit"
	"github.com/bryanl/doit/do"
	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

var (
	testKey     = do.SSHKey{Key: &godo.Key{ID: 1, Fingerprint: "fingerprint"}}
	testKeyList = do.SSHKeys{testKey}
)

func TestSSHKeysCommand(t *testing.T) {
	cmd := SSHKeys()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "create", "delete", "get", "import", "list", "update")
}

func TestKeysList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.keys.On("List").Return(testKeyList, nil)

		err := RunKeyList(config)
		assert.NoError(t, err)
	})
}

func TestKeysGetByID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.keys.On("Get", "1").Return(&testKey, nil)

		config.Args = append(config.Args, "1")

		err := RunKeyGet(config)
		assert.NoError(t, err)
	})
}

func TestKeysGetByFingerprint(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.keys.On("Get", testKey.Fingerprint).Return(&testKey, nil)

		config.Args = append(config.Args, testKey.Fingerprint)

		err := RunKeyGet(config)
		assert.NoError(t, err)
	})
}

func TestKeysCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		kcr := &godo.KeyCreateRequest{Name: "the key", PublicKey: "fingerprint"}
		tm.keys.On("Create", kcr).Return(&testKey, nil)

		config.Args = append(config.Args, "the key")

		config.Doit.Set(config.NS, doit.ArgKeyPublicKey, "fingerprint")

		err := RunKeyCreate(config)
		assert.NoError(t, err)
	})
}

func TestKeysDeleteByID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.keys.On("Delete", "1").Return(nil)

		config.Args = append(config.Args, "1")

		err := RunKeyDelete(config)
		assert.NoError(t, err)
	})
}

func TestKeysDeleteByFingerprint(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.keys.On("Delete", "fingerprint").Return(nil)

		config.Args = append(config.Args, "fingerprint")

		err := RunKeyDelete(config)
		assert.NoError(t, err)
	})

}

func TestKeysUpdateByID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		kur := &godo.KeyUpdateRequest{Name: "the key"}
		tm.keys.On("Update", "1", kur).Return(&testKey, nil)

		config.Args = append(config.Args, "1")

		config.Doit.Set(config.NS, doit.ArgKeyName, "the key")

		err := RunKeyUpdate(config)
		assert.NoError(t, err)
	})

}

func TestKeysUpdateByFingerprint(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		kur := &godo.KeyUpdateRequest{Name: "the key"}
		tm.keys.On("Update", "fingerprint", kur).Return(&testKey, nil)

		config.Args = append(config.Args, "fingerprint")

		config.Doit.Set(config.NS, doit.ArgKeyName, "the key")

		err := RunKeyUpdate(config)
		assert.NoError(t, err)
	})

}

func TestSSHPublicKeyImportWithName(t *testing.T) {
	pubkey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCn6eZ8ve0ha04rPRZuoPXK1AQ/h21qslWCzoDcOciXn5OcyafkZw+31k/afaBTeW62D8fXd8e/1xWbFfp/2GqmslYpNCTPrtpNhsE8I0yKjJ8FxX9FfsCOu/Sv83dWgSpiT7pNWVKarZjW9KdKKRQljq1i+H5pX3r5Q9I1v+66mYTe7qsKGas9KWy0vkGoNSqmTCl+d+Y0286chtqBqBjSCUCI8oLKPnJB86Lj344tFGmzDIsJKXMVHTL0dF8n3u6iWN4qiRU+JvkoIkI3v0JvyZXxhR2uPIS1yUAY2GC+2O5mfxydJQzBdtag5Uw8Y7H5yYR1gar/h16bAy5XzRvp testkey"
	path := filepath.Join(os.TempDir(), "key.pub")
	err := ioutil.WriteFile(path, []byte(pubkey), 0600)
	assert.NoError(t, err)
	defer os.Remove(path)

	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		kcr := &godo.KeyCreateRequest{Name: "custom", PublicKey: pubkey}
		tm.keys.On("Create", kcr).Return(&testKey, nil)

		config.Args = append(config.Args, "custom")

		config.Doit.Set(config.NS, doit.ArgKeyPublicKeyFile, path)

		err := RunKeyImport(config)
		assert.NoError(t, err)
	})
}
