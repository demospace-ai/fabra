package crypto

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"go.fabra.io/server/common/application"
	"go.fabra.io/server/common/errors"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const CONNECTION_KEY = "projects/fabra-344902/locations/global/keyRings/data-connection-keyring/cryptoKeys/data-connection-key"
const API_KEY_KEY = "projects/fabra-344902/locations/global/keyRings/api-key-keyring/cryptoKeys/api-key-key"
const WEBHOOK_SIGNING_KEY_KEY = "projects/fabra-344902/locations/global/keyRings/webhook-verification-key-keyring/cryptoKeys/webhook-verification-key-key"
const END_CUSTOMER_API_KEY_KEY = "projects/fabra-344902/locations/global/keyRings/end-customer-api-key-keyring/cryptoKeys/end-customer-api-key-key"

type CryptoService interface {
	DecryptConnectionCredentials(encryptedCredentials string) (*string, error)
	EncryptConnectionCredentials(credentials string) (*string, error)
	DecryptApiKey(encryptedApiKey string) (*string, error)
	EncryptApiKey(apiKey string) (*string, error)
	DecryptWebhookSigningKey(encryptedWebhookSigningKey string) (*string, error)
	EncryptWebhookSigningKey(webhookSigningKey string) (*string, error)
	DecryptEndCustomerApiKey(encryptedEndCustomerApiKey string) (*string, error)
	EncryptEndCustomerApiKey(endCustomerApiKey string) (*string, error)
}

type CryptoServiceImpl struct {
}

func NewCryptoService() CryptoService {
	return CryptoServiceImpl{}
}

func HashString(input string) string {
	h := sha256.Sum256([]byte(input))
	return base64.StdEncoding.EncodeToString(h[:])
}

func GenerateSigningKey() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(randomBytes)
}

func encrypt(keyName string, plaintextString string) (*string, error) {
	// TODO: encrypt with local keys here
	// don't encrypt in dev
	if !application.IsProd() {
		hexEncoded := hex.EncodeToString([]byte(plaintextString))
		return &hexEncoded, nil
	}

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.encrypt) failed to create kms client")
	}
	defer client.Close()

	plaintext := []byte(plaintextString)

	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	plaintextCRC32C := crc32c(plaintext)

	req := &kmspb.EncryptRequest{
		Name:            keyName,
		Plaintext:       plaintext,
		PlaintextCrc32C: wrapperspb.Int64(int64(plaintextCRC32C)),
	}

	result, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.encrypt) failed to encrypt")
	}

	if !result.VerifiedPlaintextCrc32C {
		return nil, errors.Newf("(crypto.encrypt) request corrupted in-transit")
	}
	if int64(crc32c(result.Ciphertext)) != result.CiphertextCrc32C.Value {
		return nil, errors.Newf("(crypto.encrypt) response corrupted in-transit")
	}

	ciphertext := hex.EncodeToString(result.Ciphertext)
	return &ciphertext, nil
}

func decrypt(keyName string, ciphertextString string) (*string, error) {
	ciphertext, err := hex.DecodeString(ciphertextString)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.decrypt)")
	}

	// TODO: decrypt with local keys here
	// don't encrypt in dev
	if !application.IsProd() {
		ciphertextStr := string(ciphertext)
		return &ciphertextStr, nil
	}

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.decrypt) failed to create kms client")
	}
	defer client.Close()

	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	ciphertextCRC32C := crc32c(ciphertext)

	req := &kmspb.DecryptRequest{
		Name:             keyName,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	result, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.decrypt) failed to decrypt ciphertext")
	}

	if int64(crc32c(result.Plaintext)) != result.PlaintextCrc32C.Value {
		return nil, errors.Newf("(crypto.decrypt) response corrupted in-transit")
	}

	plaintext := string(result.Plaintext)
	return &plaintext, nil
}

func (cs CryptoServiceImpl) DecryptConnectionCredentials(encryptedCredentials string) (*string, error) {
	credentials, err := decrypt(CONNECTION_KEY, encryptedCredentials)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.DecryptConnectionCredentials)")
	}

	return credentials, nil
}

func (cs CryptoServiceImpl) EncryptConnectionCredentials(credentials string) (*string, error) {
	encryptedCredentials, err := encrypt(CONNECTION_KEY, credentials)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.EncryptConnectionCredentials)")
	}

	return encryptedCredentials, nil
}

func (cs CryptoServiceImpl) DecryptApiKey(encryptedApiKey string) (*string, error) {
	apiKey, err := decrypt(API_KEY_KEY, encryptedApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.DecryptApiKey)")
	}

	return apiKey, nil
}

func (cs CryptoServiceImpl) EncryptApiKey(apiKey string) (*string, error) {
	encryptedApiKey, err := encrypt(API_KEY_KEY, apiKey)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.EncryptApiKey)")
	}

	return encryptedApiKey, nil
}

func (cs CryptoServiceImpl) DecryptWebhookSigningKey(encryptedWebhookSigningKey string) (*string, error) {
	webhookSigningKey, err := decrypt(WEBHOOK_SIGNING_KEY_KEY, encryptedWebhookSigningKey)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.DecryptWebhookSigningKey)")
	}

	return webhookSigningKey, nil
}

func (cs CryptoServiceImpl) EncryptWebhookSigningKey(webhookSigningKey string) (*string, error) {
	encryptedWebhookSigningKey, err := encrypt(WEBHOOK_SIGNING_KEY_KEY, webhookSigningKey)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.EncryptWebhookSigningKey)")
	}

	return encryptedWebhookSigningKey, nil
}

func (cs CryptoServiceImpl) DecryptEndCustomerApiKey(encryptedEndCustomerApiKey string) (*string, error) {
	endCustomerApiKey, err := decrypt(END_CUSTOMER_API_KEY_KEY, encryptedEndCustomerApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.DecryptEndCustomerApiKey)")
	}

	return endCustomerApiKey, nil
}

func (cs CryptoServiceImpl) EncryptEndCustomerApiKey(endCustomerApi string) (*string, error) {
	encryptedEndCustomerApiKey, err := encrypt(END_CUSTOMER_API_KEY_KEY, endCustomerApi)
	if err != nil {
		return nil, errors.Wrap(err, "(crypto.EncryptEndCustomerApiKey)")
	}

	return encryptedEndCustomerApiKey, nil
}
