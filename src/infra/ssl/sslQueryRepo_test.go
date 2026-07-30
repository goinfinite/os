package sslInfra

import (
	"os"
	"path/filepath"
	"testing"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

func TestSslQueryRepoSslPairFactory(t *testing.T) {
	t.Run("MultiDotHostnameDoesNotIncludeCertExtension", func(t *testing.T) {
		hostname, err := tkValueObject.NewFqdn("test.example.com")
		if err != nil {
			t.Fatalf("NewFqdnFailed: %s", err.Error())
		}

		tempDir := t.TempDir()

		writeSelfSignedPairFiles := func() {
			synth := tkInfra.Synthesizer{}
			certContent, keyContent, synthErr :=
				synth.SelfSignedCertificatePairPemFactory(&hostname, nil)
			if synthErr != nil {
				t.Fatalf(
					"SelfSignedCertificatePairPemFactoryFailed: %s",
					synthErr.Error(),
				)
			}

			certFilePath := filepath.Join(tempDir, hostname.String()+".crt")
			writeErr := os.WriteFile(
				certFilePath, []byte(certContent), 0600,
			)
			if writeErr != nil {
				t.Fatalf("WriteCertFailed: %s", writeErr.Error())
			}

			keyFilePath := filepath.Join(tempDir, hostname.String()+".key")
			writeErr = os.WriteFile(
				keyFilePath, []byte(keyContent), 0600,
			)
			if writeErr != nil {
				t.Fatalf("WriteKeyFailed: %s", writeErr.Error())
			}
		}
		writeSelfSignedPairFiles()

		crtFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(
			filepath.Join(tempDir, "test.example.com.crt"), false,
		)
		if err != nil {
			t.Fatalf("NewUnixAbsoluteFilePathFailed: %s", err.Error())
		}

		repo := NewSslQueryRepo()
		sslPairEntity, err := repo.sslPairFactory(crtFilePath)
		if err != nil {
			t.Fatalf("SslPairFactoryFailed: %s", err.Error())
		}

		if sslPairEntity.VirtualHostHostname.String() != "test.example.com" {
			t.Errorf(
				"UnexpectedVirtualHostHostname: '%s' expected 'test.example.com'",
				sslPairEntity.VirtualHostHostname.String(),
			)
		}
	})
}
