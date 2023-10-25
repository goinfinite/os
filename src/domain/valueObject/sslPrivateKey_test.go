package valueObject

import "testing"

func TestSslPrivateKey(t *testing.T) {
	t.Run("ValidSslPrivateKey", func(t *testing.T) {
		validSslPrivateKey := `
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAziJ8BEmVq/cSznb7aYRL5YBJjhMZVxC/jqT2Q/LKWFjX38Er
LHT3khhdlKyrh/7AgfN+Us1Q7/eHq3PKX7Z0lgk+9LssNnaH67bj9lqJcILlToc5
UrhZLHe2Q3xUlfyveDoheepcbBiqO7xzuUNl53KpT4FlF3DBO94wMNNqOjeNslHv
lfpc/gJlZ5IBuxGG1+xjA75bGqDnFqyjEyUxrNxJyM70NAL4J+3rScdvXpdMKbMn
IqC2s0mrK1iuPL4tryOG5/dES5BPJdRIrEAC/G4Kf6h8xD/QI9zzPmiZ4UJqb0A2
TcmYexc40BnkgHO7XWr1zP20oZvSS46C/v3kAwIDAQABAoIBAQCd02Vk2vpP2jJ6
BjtkhLiflWO79f+W2+nuy3sKd2BZ2Fwgo4Ps2/mZ0DIGXVZQH8tBNC9qMm1f7gPg
UB2Ivufw4E9ljdHCOWrEHRnZS2Sj0nTDdWF8Zk1QcLAKZ61T0U6AHPH4qGnvEct1
RUrNdD8XwIDFsOq30csBjZMULyrMOtC5B39mi46cWiQVjoabiOm7IDZFUFEq9Oy9
ZGcyUWRce60NVcUgYiH1OTY3JNSdAIWXyeaViGNvIdQFt9XLnWkSCF2OZ+WzANy1
2GePJdG6Jnn2XYsvpKy2mQXp+ULwejPN/KTcuIKPAL29B/wCOtZAXCiE6yib5g9a
98YlWcjBAoGBAO2NWG7hOFyveijjJNEq9MyhhALu/o3xM8HwPQJOL1moWt4uySUc
nbnl14YiuZjpoGbuL5143NXN2sJA+XnTBq2pRAPBJS1OLLkY1W+BVEOAc72TRV9v
egWOHomgUqyInfRTmutnFCVyaAOURRpiXkPBMRLfN9WemVBBbuOAxO/zAoGBAN4k
i4FEXqHFKp2OR7iRQenP/mSQGTQcQ1zUeWun02nEfMzosFcyioFZwQi3NMvtdC7u
meHDuoGwLH+YcoM3/5NKUn1Kccl3x3tIbRy9GhavX2NBQ4v0DX2h7acpN4soWlHs
v07AUMiYHvxEi+t2R7UptUufiwIwSlCmBQwGGE+xAoGBAOYu4FIQypyFLMoRz8se
5Laki1aMXv0LjCuQro1dVWR7ThGdJCth3zQTExRW8aDKQTN7+YeNZe+G2UMB0rvJ
T99W9SDuNyf/aDazaZ3yo8QE5CH+YmpnisV3QP/66iFlACmQGb2g1FS01zUgpxU5
3D2rJfIzedb1J3os7VZloG8hAoGBANFxJ07DhW2Elf9izGBKJBksj6+E5R5qn2CA
u9Iys3N/XCNeKBSuhEQcuZFcGp1Czk4JjHB9t/Tag7nxo9XwEDlw04FplQrcsemc
ibOU32oQAyFzwRnNCoMvDwCSLdo4O6AOVPkM/Z2DP4OdpUZliIpYPqSEUe3IVejf
/tYtUPKhAoGAQBzkCzCunVBtSuVLPT5/9NVVCj1nfxmXdLq1+0bD5q6m2XiW1G7T
YvdMz2Th5XVtTfNHhLIKpMyrq7sstb6lsQPKNKSpBuyHa8oooWHrkxf7VVBrN7v8
en/D01Fd9hVhXyGKaOk/nEDxB8fTgQ1dE9JhxUiqHN4Po1ktNG/P8aU=
-----END RSA PRIVATE KEY-----`

		_, err := NewSslPrivateKey(validSslPrivateKey)
		if err != nil {
			t.Errorf("Expected no error for dummy SSL private key, got %v", err)
		}
	})

	t.Run("InvalidSslPrivateKey", func(t *testing.T) {
		invalidSslPrivateKeys := []string{
			`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkj
-----END PUBLIC KEY-----`,
			`MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==`,
		}
		for sslPrivateKeyIndex, sslPrivateKey := range invalidSslPrivateKeys {
			_, err := NewSslPrivateKey(sslPrivateKey)
			if err == nil {
				t.Errorf("Expected error for '%v' SSL private key, got nil", sslPrivateKeyIndex)
			}
		}
	})
}
