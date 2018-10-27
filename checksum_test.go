package appcast

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecksum(t *testing.T) {
	c := NewChecksum(SHA256, []byte("test"))
	assert.IsType(t, Checksum{}, *c)
	assert.Equal(t, SHA256, c.algorithm)
	assert.Equal(t, []byte("test"), c.source)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", c.String())
}

func TestChecksum_generate(t *testing.T) {
	testCases := map[string][]string{
		"github/default.xml": {
			"c28ff87daf2c02471fd2c836b7ed3776d927a8febbb6b8961daf64ce332f6185",
			"572802c8d0cae5435461d73844764463",
		},
		"github/invalid_pubdate.xml": {
			"52f87bba760a4e5f8ee418cdbc3806853d79ad10d3f961e5c54d1f5abf09b24b",
			"1aeca62fcdb36aa5ed3c18efdbcc9c02",
		},
		"github/invalid_version.xml": {
			"7375a6cbee6f9369bd8e4ecbda347889a0272b8dd8a5eb473c1dec9dfa753392",
			"ca0ee1fef654c37bb1e8789ad004bf09",
		},
		"sourceforge/default.xml": {
			"d4afcf95e193a46b7decca76786731c015ee0954b276e4c02a37fa2661a6a5d0",
			"76848b058151cae70fcf7d3838329517",
		},
		"sourceforge/empty.xml": {
			"569cb5c8fa66b2bae66e7c0d45e6fbbeb06a5f965fc7e6884ff45aab4f17b407",
			"f78eefccfcf70937a004b94bd063682b",
		},
		"sourceforge/invalid_pubdate.xml": {
			"160885aaaa2f694b5306e91ea20d08ef514f424e51704947c9f07fffec787cf6",
			"c39e1ffc7bbe1e86fe252052269fb766",
		},
		"sourceforge/invalid_version.xml": {
			"ad841a02d68c60589136f1f01d000b7988989c187da3ffabbf9d89832a84a6f1",
			"7b10e5b85c41d8601ac939fbf50a1da5",
		},
		"sourceforge/single.xml": {
			"5384ed38515985f60f990c125f1cceed0261c2c5c2b85181ebd4214a7bc709de",
			"d81ac48573ac8a08eba66b44104eac7e",
		},
		"sparkle/attributes_as_elements.xml": {
			"d59d258ce0b06d4c6216f6589aefb36e2bd37fbd647f175741cc248021e0e8b4",
			"f2a145079bd9f012de062d60e1bb3190",
		},
		"sparkle/default_asc.xml": {
			"9f8d8eb4c8acfdd53e3084fe5f59aa679bf141afc0c3887141cd2bdfe1427b41",
			"dad52a6efffd480f2d669ca0795dbd99",
		},
		"sparkle/default.xml": {
			"0cb017e2dfd65e07b54580ca8d4eedbfcf6cef5824bcd9539a64afb72fa9ce8c",
			"21448d1059f783c979967c116b255d43",
		},
		"sparkle/incorrect_namespace.xml": {
			"ff464014dc6a2f6868aca7c3b42521930f791de5fc993d1cc19d747598bcd760",
			"3e46208ce5fb1947aef2c9c7917aa770",
		},
		"sparkle/invalid_pubdate.xml": {
			"9a59f9d0ccd08b317cf784656f6a5bd0e5a1868103ec56d3364baec175dd0da1",
			"c108c194a53216044fbae679a8c0bc76",
		},
		"sparkle/invalid_version.xml": {
			"65d754f5bd04cfad33d415a3605297069127e14705c14b8127a626935229b198",
			"2b63f313288387799a9962e81e35ba5f",
		},
		"sparkle/multiple_enclosure.xml": {
			"b3b1304739b58126eef8386f134a82d5c71d7b83a076ec732b2fb133734524f3",
			"9288bd846c263f0d35eb2fa5b7db986f",
		},
		"sparkle/no_releases.xml": {
			"befd99d96be280ca7226c58ef1400309905ad20d2723e69e829cf050e802afcf",
			"6b2f1f5e0cea6005e5410c1d76cab0a3",
		},
		"sparkle/only_version.xml": {
			"ee5a775fec4d7b95843e284bff6f35f7df30d76af2d1d7c26fc02f735383ef7f",
			"d907f5f8366bc27f99967dd66877b15c",
		},
		"sparkle/prerelease.xml": {
			"8e44fccf005ad4720bcc75b9afffb035befade81bdf9f587984c26842dd7c759",
			"0c2dc95224959f8e4045f1b3983171db",
		},
		"sparkle/single.xml": {
			"c59ec641579c6bad98017db7e1076a2997cdef7fff315323dd7f0cabed638d50",
			"714c612b683650f52d313ad3bf4070f3",
		},
		"sparkle/with_comments.xml": {
			"354fc20caba7cb67a24378c309ee696ebb16a91db38cb2bdd53b700f258e9a82",
			"83995df99eef503369c94d40afa13890",
		},
		"sparkle/without_namespaces.xml": {
			"888494294fc74990e4354689a02e50ff425cfcbd498162fdffd5b3d1cd096fa1",
			"e85bac3581dac054e2d9ac42c59bf07f",
		},
		"unknown.xml": {
			"c29665078d79a8e67b37b46a51f2a34c6092719833ccddfdda6109fd8f28043c",
			"2340f9a888f7305f4636d4f70d3471b1",
		},
	}

	for filename, checkpoints := range testCases {
		content := getTestdata(filename)

		// SHA256
		c := &Checksum{SHA256, content, nil}
		c.generate()
		assert.Equal(t, checkpoints[0], c.String(), fmt.Sprintf("Checksum doesn't match (SHA256): %s", filename))

		// MD5
		c = &Checksum{MD5, content, nil}
		c.generate()
		assert.Equal(t, checkpoints[1], c.String(), fmt.Sprintf("Checksum doesn't match (MD5): %s", filename))
	}
}

func TestChecksum_Algorithm(t *testing.T) {
	c := NewChecksum(SHA256, []byte("test"))
	assert.Equal(t, SHA256, c.Algorithm())
}

func TestChecksum_Source(t *testing.T) {
	src := []byte("test")
	c := NewChecksum(SHA256, src)
	assert.Equal(t, src, c.Source())
}

func TestChecksum_Result(t *testing.T) {
	result, _ := hex.DecodeString("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
	c := NewChecksum(SHA256, []byte("test"))
	assert.Equal(t, result, c.Result())
}

func TestChecksumAlgorithm_String(t *testing.T) {
	assert.Equal(t, "SHA256", SHA256.String())
	assert.Equal(t, "MD5", MD5.String())
}

func TestChecksum_String(t *testing.T) {
	c := NewChecksum(SHA256, []byte("test"))
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", c.String())
}
