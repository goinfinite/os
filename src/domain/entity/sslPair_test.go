package entity

import (
	"testing"

	"github.com/speedianet/os/src/domain/valueObject"
)

func TestSslPair(t *testing.T) {
	t.Run("SelfSignedSslPair", func(t *testing.T) {
		certContentStr := `-----BEGIN CERTIFICATE-----
MIIDlTCCAn2gAwIBAgIIUWRFUJIkWWAwDQYJKoZIhvcNAQEFBQAwNDEZMBcGA1UE
AxMQdGVzdC5leGFtcGxlLmNvbTEXMBUGA1UEChMOc2VsZnNpZ25lZC5vcmcwHhcN
MjQwNjE5MTkzNzU0WhcNMzQwNjE5MTkzNzU0WjA0MRkwFwYDVQQDExB0ZXN0LmV4
YW1wbGUuY29tMRcwFQYDVQQKEw5zZWxmc2lnbmVkLm9yZzCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAMi+bP0rebHu9mf5+EW729pSE+RD9HBO1x7Mnr02
Twc4GSTVGoC3/DxJnIpE7bDs190uJBymYzqvQjQ99wYTF1UbR3kiViRgfwp6RepI
ZEwDLwc3xhrcqmfC5MpkDyq4PcucsRCVA7h5dvjzUHigVxNTD23KIuu4Yss6KHuc
a8SJ+DCzYRFEqFxTjR4IYfuOwUt8QfJUpFAEZjtU4OzTugphM2elg9Pi1soMArV+
r3nJ7knoe0J15t/GpSWyiUwzYVYarEOALAeBtpVmwZ7kCPZT1fJe/Wc47qbtQiwG
VePzGUUdzabJSsUh9vuj0J8fswiXq3CCgwly6V3mGN45nocCAwEAAaOBqjCBpzAM
BgNVHRMEBTADAQH/MAsGA1UdDwQEAwIC9DA7BgNVHSUENDAyBggrBgEFBQcDAQYI
KwYBBQUHAwIGCCsGAQUFBwMDBggrBgEFBQcDBAYIKwYBBQUHAwgwEQYJYIZIAYb4
QgEBBAQDAgD3MBsGA1UdEQQUMBKCEHRlc3QuZXhhbXBsZS5jb20wHQYDVR0OBBYE
FHvu8BMe2TZKxuzA9VSthSvZd7fbMA0GCSqGSIb3DQEBBQUAA4IBAQBeZwQ6GZkA
tKNk14tyDROaq/Ngu+vpnKzo+pYya5bMsxJcDowR8Lh+UQiqf4S+iCEtphMm3F4T
nRR/Jp6weBlozbaVVxutZBJCMfzLrzKTI3B3ndxHyljFES4/syZD83QTHxIH7RLo
hdqVmCQYskmZej2viMae8Ca7GBZPCcuKplhga4KEDI5DI20Ojj8Tj/EKX9CSMzNo
P5MLgmkcAlfjeSVmbeRT2gtypaRw5zYUm96Yt/yMdhkkLV/Uki6wXlPQk6seqZzA
Jv8gduA0cyBXAG1Ba+aKTl44TPlHfDnrOWuK3aUi3aHMOFimAvTWGuklD2ylhTT1
eicXskJoC+KU
-----END CERTIFICATE-----`
		certificateContent, err := valueObject.NewSslCertificateContent(certContentStr)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		certificate, err := NewSslCertificate(certificateContent)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		privateKeyStr := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAyL5s/St5se72Z/n4Rbvb2lIT5EP0cE7XHsyevTZPBzgZJNUa
gLf8PEmcikTtsOzX3S4kHKZjOq9CND33BhMXVRtHeSJWJGB/CnpF6khkTAMvBzfG
GtyqZ8LkymQPKrg9y5yxEJUDuHl2+PNQeKBXE1MPbcoi67hiyzooe5xrxIn4MLNh
EUSoXFONHghh+47BS3xB8lSkUARmO1Tg7NO6CmEzZ6WD0+LWygwCtX6vecnuSeh7
QnXm38alJbKJTDNhVhqsQ4AsB4G2lWbBnuQI9lPV8l79Zzjupu1CLAZV4/MZRR3N
pslKxSH2+6PQnx+zCJercIKDCXLpXeYY3jmehwIDAQABAoIBAEfSNZ1IqjQaimc9
/HE6mpicSAihtXlfA8E9tUd+AD1VeU1/vwkxilmpfovLyHzF6B92rC3h69upq5aU
WuZ9+xmUdnhk7Av8yEcf4xbEyrmVZASBlGu06nTQOloc/X4rx9Qq3gDQR7H/Jy0/
pFlcCHtd+sWtjdvnLtWGG8jJ+JaqoNNACNV1EYG4IWXwrtcci4F4apdRRE/16+rC
IWnAcYd9ueX33GHkenAMCAQmSbwp0XF42AC1urM8TvOe2q/tftZZaDSprVFiaeKH
CZFdAXLuzmbDXqj+CiFd7OUB7ZRrGcAHvHoNCHm9cWi5/dkrbekYJzewaFZqJJ/s
bIIN2zECgYEA/5V2zaobiLZ5O/NVE+ZwkOZILY8zZr7QQtNI09q9I12QjQYJ5cJ4
ECcVzXDWGx0CzOlEf0LH24SbeHGa8+g7S6Fnp5HZoEsO3VFwoZTKKJMhu5ybmXKr
yQfIM5d6/MYOa9Awk0Fbm9cRyAYxRtZKd5BlR4YbtfNiZbIC3IsqUHcCgYEAyRIa
PkJMC5adfzwMWWa4amF1m93ByznTBLxcYkatJ9/WanSmycft88KrAQ/P4VLwrX2X
jI12y3V4VIRLkl5cViVIetTcyAC7l6pOxHXVbv5kdWKOYmVAFLMj9lk09sCtxkcE
P6fSHMG0SPEZARhAwecXIU0WjEotJt+DTijFNnECgYEA+U4bF9xhhUaxFUhzabjz
fnQSXdZ8hjGExlqAhJ6uteuTj+wfBW5fXSoy+zWgs8vlqmmz9gr3FmrQmHkAdADI
ripgCLWdOd1dP4csPYD8fP2f/vhxUwnnBW5A3Apb3mt3L7VhXJJ5QJdWce2QbY+k
DeLc2Bq5tw8UoSw13FknSlsCgYA2DxvnKUPwyanGj4pybt+eGl3YbiKwVPebCll8
QqxDUDcBoCNHlO0w4GHBg1LMrdPvkRixvUb3JLoZXwhCbgQ9VQDLpXdGfovxFuTe
hR/BG7w+oyTM55P2/MLqdMl8ngkaifVmd+RRvvKNueSTGsYuW8coOOWbCkZhcS6I
UQXUwQKBgQCeJ0lSt3fdgE4fPVcd0kggWMD/H6J9W3IKsxqegkKRHv/6BwJ4ogNk
6pbF0picFayT99XyKWR/Oz5+fxI42ZsYLUjnxPwcsocNByty0blOQuo5hb7lKTLk
ZrBUC3x7Z1Ex9qgz3p/Y/WxHDRrHDuuByGazDVhyDyXECWxpDYEK3g==
-----END RSA PRIVATE KEY-----`
		privateKey, err := valueObject.NewSslPrivateKey(privateKeyStr)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		pairId, err := valueObject.NewSslIdFromSslCertificateContent(certificateContent)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		chainCertificates := []SslCertificate{}

		testHostname, err := valueObject.NewFqdn("test.example.com")
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		vhostHostnames := []valueObject.Fqdn{testHostname}

		sslPair := NewSslPair(pairId, vhostHostnames, certificate, privateKey, chainCertificates)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		if sslPair.IsPubliclyTrusted() {
			t.Errorf("SelfSignedSslPairShouldNotBePubliclyTrusted")
		}
	})

	// This test is going to fail when the certificate expires. Make sure you update
	// the certificate content when you need to run this test.
	t.Run("PubliclyTrustedSslPair", func(t *testing.T) {
		certContentStr := `-----BEGIN CERTIFICATE-----
MIIE9jCCA96gAwIBAgISBGxRd4HZZLll27lcr/pbpTCnMA0GCSqGSIb3DQEBCwUA
MDMxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQwwCgYDVQQD
EwNSMTAwHhcNMjQwNjE5MTY0MDE2WhcNMjQwOTE3MTY0MDE1WjAeMRwwGgYDVQQD
ExNvcy5kZW1vLnNwZWVkaWEubmV0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAw/7d/qAjafhyonbaEWXnmEUks9mPrqeWiodrJ4eg4i6CWsjt6P3hFHp5
FToqQn2Pg6fB/gF44RKz+Sa2ulPLqGOafo/8YKsiFhaLLrCiatXng/3v662J6IzM
TqTzyLy88xLyiyw1bCKbvqoM+NH1MJNmiActdrywWg1/lwKvAeNX1xJv0jkdVmcd
xnuwhCdd06teaDxNZu0vQ3Imvni7+vmIxkIO+khJAgPBP8l8CbP5owDtDzqvfGb3
46wAPIhkzAoncgISRbYbtd2/xkjI8ym+AV8MLFfs11qLC2ScGE7l7yc51gOI0pXb
NKrCnusIMqLfY2dToK3rR1f1N1zk7wIDAQABo4ICFzCCAhMwDgYDVR0PAQH/BAQD
AgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAA
MB0GA1UdDgQWBBTytxAmylz9nXif/jkKCYTutmmOgDAfBgNVHSMEGDAWgBS7vMNH
peS8qcbDpHIMEI2iNeHI6DBXBggrBgEFBQcBAQRLMEkwIgYIKwYBBQUHMAGGFmh0
dHA6Ly9yMTAuby5sZW5jci5vcmcwIwYIKwYBBQUHMAKGF2h0dHA6Ly9yMTAuaS5s
ZW5jci5vcmcvMB4GA1UdEQQXMBWCE29zLmRlbW8uc3BlZWRpYS5uZXQwEwYDVR0g
BAwwCjAIBgZngQwBAgEwggEEBgorBgEEAdZ5AgQCBIH1BIHyAPAAdgA/F0tP1yJH
WJQdZRyEvg0S7ZA3fx+FauvBvyiF7PhkbgAAAZAxlcKDAAAEAwBHMEUCIQDLS9iq
XslKIWfQ10N8Y+i26i30J6RCwd87CdfOwq7nNwIgd/O+61no29P050HCJyZ0ON5z
J98ZOyq8kroG9euqX9UAdgB2/4g/Crb7lVHCYcz1h7o0tKTNuyncaEIKn+ZnTFo6
dAAAAZAxlcLDAAAEAwBHMEUCIQCYsTwdkGEDc4YNm5VLZcKhLxwm35LmvAWRNGLj
3NqO/AIgMJ5wTQOFcxeIgE0H4Nhe+s7lwRxszRj5L9ak3/GJYrgwDQYJKoZIhvcN
AQELBQADggEBADHAn+FY/fRRo17+N1EDQy4kWBUirRhvhe3liFzsjaLMFG2r8dgr
cbH4ZuImRDhu0ma05qA0JgCzwKnxRuE4Qox+vqfpzCJlADMvCAQHM9qELgOuS7YK
h7GJRnOuUDFCHz0SaNxzIt1AyplfVsmUsFu0+dvdZVZs0W3O4pi+umDEW1DQZ67V
4RJ4mAt5FThn7D37tpiMV97KtbS0k/jsBlsCLWDXaQG3y1nHmxg2gppN+OTk7AzB
S8zlukal9/D+kEwXITuf9Fb1ctjnxUvwaVoZYHwtgnvLEeoYqewuPPv7w0hk/x6J
ZQGEW5pclcWegbW1lvS2qSdtKwgW89ETEPI=
-----END CERTIFICATE-----`
		certificateContent, err := valueObject.NewSslCertificateContent(certContentStr)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		certificate, err := NewSslCertificate(certificateContent)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		privateKeyStr := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDD/t3+oCNp+HKi
dtoRZeeYRSSz2Y+up5aKh2snh6DiLoJayO3o/eEUenkVOipCfY+Dp8H+AXjhErP5
Jra6U8uoY5p+j/xgqyIWFosusKJq1eeD/e/rrYnojMxOpPPIvLzzEvKLLDVsIpu+
qgz40fUwk2aIBy12vLBaDX+XAq8B41fXEm/SOR1WZx3Ge7CEJ13Tq15oPE1m7S9D
cia+eLv6+YjGQg76SEkCA8E/yXwJs/mjAO0POq98ZvfjrAA8iGTMCidyAhJFthu1
3b/GSMjzKb4BXwwsV+zXWosLZJwYTuXvJznWA4jSlds0qsKe6wgyot9jZ1OgretH
V/U3XOTvAgMBAAECggEBALjIpcvWdw0F7C44F8inZv4s0jmuOMTVxBy/J8uRF6Gn
b8bgAU3Vbku8XEQjHoypjJD3rPIpuSmaVIvmcAij0DLmFVaVscACGJTylC5k4fwP
x8Ktu3Fbn9XcSRMseZscNpiFmJ6WA5f72RKdiLVeXeh6UASXn8l+hFWivFRHd3Ay
mgXxPxVSXsXpQg8WWRpalxA6AT9TuUEkgSuuH+EiDpydqqxn7W8jDhUgMAFdA7pF
/CbUpZ2ojkMWtpRHZKpWViPm7/sqoo0CjtUZurqxCqNT3FazAukK/l7U5gBDteRU
2DLq+6dwqiKXfA77a1cKuCqwEQoz/r9bkhgDcmgTZ8ECgYEA+1DJ9vdcVdBWyq+H
eygcvZ8LDt+gfnuN5RqxoXdiqnCmc1gjVx6toPh7QCnVOQZxsOLZEj+P/L50bgoL
jD92CZvR8qloRUbGrxmdR233sPFiQ+ULFMD+uDXOS32SU0jhg8dYA91Vv17hexcB
mUwtTAHgZURQbAMT32cB+CxAcccCgYEAx6YbG6wJymw6BZgFbUG1ZxFFooITLgKz
DjsvFF7LOhD6NO0sQCPLD3vgSSErcdV2EGqSCOg8OWn8Tzx97ZN1GE7O+CsMhvGw
YcosOAYzEy7H049cywmC57nHwfoiVrK2PDiq1uZZF+kh+hmegxZ9kYr25d5NV6rA
PZztjaBf85kCgYBFo21Tcde0L7bdEyaHieXs5VU7GdxvL+1xvqPaCirc77ov3Axu
56FVKYV9khnzY0W9rh5YYCSV9HBuzXnFsxASOYtDoo2yJJqJip96W453CWwhRCZ9
6byrbr1rTbBuQ5O54FMTPxGzpab1ZLqcr+8dUKfNZ9ChHXk0PmbdTeXNoQKBgEQy
FvxizY1ZXpBelyv0z/P+0FxsNgT3YxYvXSuGHcEd33mIsh7OmyQU2k3giKh/k66T
2II2Yavy6f5e2Vz3i33cHZJjkgneMLOWjXgtlfCtgBBh3f50p0RkDznRKT2YGeuE
J8b0M+aY+cQmUCDS4919LEzknGKfrr7dBb/k0iGxAoGAEj4571UqjCYpbfXxlkKo
HGhQoFN3vSgKYdGKYcHGMTjM+L++fBf18MQ5KJ2cwIZmhyVd5it0x8omam+mxmHx
TjK39rhrlvV4HPH0SUp0+LHgP83UTOdbggJwDV3mRkuNrxdqKY8slJfOurAJtk1t
cgNzy7pXTh6r4/5EvmhJ+1o=
-----END PRIVATE KEY-----`
		privateKey, err := valueObject.NewSslPrivateKey(privateKeyStr)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		pairId, err := valueObject.NewSslIdFromSslCertificateContent(certificateContent)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		chainCertContentStr := `-----BEGIN CERTIFICATE-----
MIIFBTCCAu2gAwIBAgIQS6hSk/eaL6JzBkuoBI110DANBgkqhkiG9w0BAQsFADBP
MQswCQYDVQQGEwJVUzEpMCcGA1UEChMgSW50ZXJuZXQgU2VjdXJpdHkgUmVzZWFy
Y2ggR3JvdXAxFTATBgNVBAMTDElTUkcgUm9vdCBYMTAeFw0yNDAzMTMwMDAwMDBa
Fw0yNzAzMTIyMzU5NTlaMDMxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBF
bmNyeXB0MQwwCgYDVQQDEwNSMTAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDPV+XmxFQS7bRH/sknWHZGUCiMHT6I3wWd1bUYKb3dtVq/+vbOo76vACFL
YlpaPAEvxVgD9on/jhFD68G14BQHlo9vH9fnuoE5CXVlt8KvGFs3Jijno/QHK20a
/6tYvJWuQP/py1fEtVt/eA0YYbwX51TGu0mRzW4Y0YCF7qZlNrx06rxQTOr8IfM4
FpOUurDTazgGzRYSespSdcitdrLCnF2YRVxvYXvGLe48E1KGAdlX5jgc3421H5KR
mudKHMxFqHJV8LDmowfs/acbZp4/SItxhHFYyTr6717yW0QrPHTnj7JHwQdqzZq3
DZb3EoEmUVQK7GH29/Xi8orIlQ2NAgMBAAGjgfgwgfUwDgYDVR0PAQH/BAQDAgGG
MB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATASBgNVHRMBAf8ECDAGAQH/
AgEAMB0GA1UdDgQWBBS7vMNHpeS8qcbDpHIMEI2iNeHI6DAfBgNVHSMEGDAWgBR5
tFnme7bl5AFzgAiIyBpY9umbbjAyBggrBgEFBQcBAQQmMCQwIgYIKwYBBQUHMAKG
Fmh0dHA6Ly94MS5pLmxlbmNyLm9yZy8wEwYDVR0gBAwwCjAIBgZngQwBAgEwJwYD
VR0fBCAwHjAcoBqgGIYWaHR0cDovL3gxLmMubGVuY3Iub3JnLzANBgkqhkiG9w0B
AQsFAAOCAgEAkrHnQTfreZ2B5s3iJeE6IOmQRJWjgVzPw139vaBw1bGWKCIL0vIo
zwzn1OZDjCQiHcFCktEJr59L9MhwTyAWsVrdAfYf+B9haxQnsHKNY67u4s5Lzzfd
u6PUzeetUK29v+PsPmI2cJkxp+iN3epi4hKu9ZzUPSwMqtCceb7qPVxEbpYxY1p9
1n5PJKBLBX9eb9LU6l8zSxPWV7bK3lG4XaMJgnT9x3ies7msFtpKK5bDtotij/l0
GaKeA97pb5uwD9KgWvaFXMIEt8jVTjLEvwRdvCn294GPDF08U8lAkIv7tghluaQh
1QnlE4SEN4LOECj8dsIGJXpGUk3aU3KkJz9icKy+aUgA+2cP21uh6NcDIS3XyfaZ
QjmDQ993ChII8SXWupQZVBiIpcWO4RqZk3lr7Bz5MUCwzDIA359e57SSq5CCkY0N
4B6Vulk7LktfwrdGNVI5BsC9qqxSwSKgRJeZ9wygIaehbHFHFhcBaMDKpiZlBHyz
rsnnlFXCb5s8HKn5LsUgGvB24L7sGNZP2CX7dhHov+YhD+jozLW2p9W4959Bz2Ei
RmqDtmiXLnzqTpXbI+suyCsohKRg6Un0RC47+cpiVwHiXZAW+cn8eiNIjqbVgXLx
KPpdzvvtTnOPlC7SQZSYmdunr3Bf9b77AiC/ZidstK36dRILKz7OA54=
-----END CERTIFICATE-----`
		chainCertificateContent, err := valueObject.NewSslCertificateContent(chainCertContentStr)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		chainCertificate, err := NewSslCertificate(chainCertificateContent)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		chainCertificates := []SslCertificate{chainCertificate}

		demoHostname, err := valueObject.NewFqdn("os.demo.speedia.net")
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		vhostHostnames := []valueObject.Fqdn{demoHostname}

		sslPair := NewSslPair(pairId, vhostHostnames, certificate, privateKey, chainCertificates)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		if !sslPair.IsPubliclyTrusted() {
			t.Errorf("PubliclyTrustedSslPairShouldBePubliclyTrusted")
		}
	})
}
