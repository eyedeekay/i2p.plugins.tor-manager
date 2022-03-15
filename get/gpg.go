package tbget

import (
	"fmt"
	"log"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
)

func Verify(keyrings, detached, target string) error {
	keyRingReader, err := os.Open(keyrings)
	if err != nil {
		return fmt.Errorf("Verify: failed to open keyrings: %s\n\t%s", err, keyrings)
	}

	signature, err := os.Open(detached)
	if err != nil {
		return fmt.Errorf("Verify: failed to open detached signature: %s\n\t%s", err, detached)
	}

	verification_target, err := os.Open(target)
	if err != nil {
		return fmt.Errorf("Verify: failed to open verification target: %s\n\t%s", err, target)
	}

	entities, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		return fmt.Errorf("Verify: failed to read keyrings: %s\n\t%s", err, keyrings)
	}
	log.Printf("Verify: %s", fmt.Sprintf("Read %d keyrings", len(entities)))
	log.Printf("Verifying: %s against %s\n", target, detached)
	log.Printf("Verify: using keyring %s\n", keyrings)
	_, err = openpgp.CheckArmoredDetachedSignature(entities, verification_target, signature, nil)
	if err != nil {
		return fmt.Errorf("Verify: failed to verify signature: %s\n\t%s\n\t%s\n\t%s", err, keyrings, detached, target)
	}

	return nil
}
