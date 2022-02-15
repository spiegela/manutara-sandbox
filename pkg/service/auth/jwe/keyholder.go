package jwe

import (
	"crypto/rand"
	"crypto/rsa"
	"sync"

	sync2 "github.com/spiegela/manutara/pkg/service/secret"

	"github.com/sirupsen/logrus"

	"gopkg.in/square/go-jose.v2"
)

// KeyHolder is responsible for generating, storing and synchronizing encryption key used for token
// generation/decryption.
type KeyHolder interface {
	// Returns encrypter instance that can be used to encrypt data.
	Encrypter() jose.Encrypter
	// Returns encryption key that can be used to decrypt data.
	Key() *rsa.PrivateKey
	// Forces refresh of encryption key synchronized with kubernetes resource (secret).
	Refresh()
}

// Implements KeyHolder interface
type rsaKeyHolder struct {
	// 256-byte random RSA key pair. Synced with a key saved in a secret.
	key  *rsa.PrivateKey
	sync sync2.Synchronizer
	mux  sync.Mutex
}

// Encrypter implements key holder interface. See KeyHolder for more information.
// Used encryption algorithms:
//    - Content encryption: AES-GCM (256)
//    - Key management: RSA-OAEP-SHA256
func (k *rsaKeyHolder) Encrypter() jose.Encrypter {
	publicKey := &k.Key().PublicKey
	encrypter, err := jose.NewEncrypter(jose.A256GCM, jose.Recipient{
		Algorithm: jose.RSA_OAEP_256,
		Key:       publicKey,
	}, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	return encrypter
}

// Key implements key holder interface. See KeyHolder for more information.
func (k *rsaKeyHolder) Key() *rsa.PrivateKey {
	k.mux.Lock()
	defer k.mux.Unlock()
	return k.key
}

// Refresh implements key holder interface. See KeyHolder for more information.
func (k *rsaKeyHolder) Refresh() {
	k.sync.Refresh()
	k.update(k.sync.Get())
}

func (k *rsaKeyHolder) init() {
	k.initEncryptionKey()
	// TODO: init key from synchronized object
	// TODO: save generated key in a secret
}

// Generates encryption key used to encrypt token payload.
func (k *rsaKeyHolder) initEncryptionKey() {
	// TODO: add mutex + synchronizer to ensure atomic updates
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logrus.Fatal(err)
	}
	k.key = privateKey
}

// NewRSAKeyHolder creates new KeyHolder instance.
func NewRSAKeyHolder() KeyHolder {
	holder := &rsaKeyHolder{}
	// TODO: add synchronizer to key  holder
	holder.init()
	return holder
}
