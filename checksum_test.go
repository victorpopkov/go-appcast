package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"encoding/hex"
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
			"8c42d7835109ff61fe85bba66a44689773e73e0d773feba699bceecefaf09359",
			"90444fe711048735501877fd54dbcbd3",
		},
		"sparkle/default_asc.xml": {
			"9f94a728eab952284b47cc52acfbbb64de71f3d38e5b643d1f3523ef84495d9f",
			"0247ff43c3df1a0c6c3f2bedf5f4be05",
		},
		"sparkle/default.xml": {
			"83c1fd76a250dd50334db793a0db5da7575fc83d292c7c58fd9d31d5bcef6566",
			"56157a2dc1cec9dc02448223e31854fa",
		},
		"sparkle/incorrect_namespace.xml": {
			"2e66ef346c49a8472bf8bf26e6e778c5b4d494723223c84c35d9f272a7792430",
			"da82e1a170325e28e4fc1ed94bacaa88",
		},
		"sparkle/invalid_pubdate.xml": {
			"e0273ccbce5a6fb6a5fe31b5edffb8173d88afa308566cf9b4373f3fed909705",
			"d98e602b718c7949a88fd41d9cc28cc8",
		},
		"sparkle/invalid_version.xml": {
			"12c7827fed4cccb5c4bc77052d2c95b03c0e4943aa49c90f9f2e98bb8ab9b799",
			"6a2a0417379a4f70272165fe053c76d0",
		},
		"sparkle/multiple_enclosure.xml": {
			"7f62916d4d80cc9a784ffa1d2211488104c4578cc2704baaff48a96b4df00961",
			"27737897524ca35a512c0ef4d9cff44a",
		},
		"sparkle/no_releases.xml": {
			"befd99d96be280ca7226c58ef1400309905ad20d2723e69e829cf050e802afcf",
			"6b2f1f5e0cea6005e5410c1d76cab0a3",
		},
		"sparkle/only_version.xml": {
			"5c3e7cf62383d4c0e10e5ec0f7afd1a5e328137101e8b6bade050812e4e7451f",
			"8adc62b60049e22985ecef4df1fd8abc",
		},
		"sparkle/prerelease.xml": {
			"56f95889fe5ddabd847adfe995304fd78dbeeefe47354c2e1c8bde0f003ecf5c",
			"9966669ccb5c80e556ded5c37130b303",
		},
		"sparkle/single.xml": {
			"ac649bebe55f84d85767072e3a1122778a04e03f56b78226bd57ab50ce9f9306",
			"aa5b165d930d81645b5c14c66bf67957",
		},
		"sparkle/with_comments.xml": {
			"283ea10e6f7f81466beb85e055940765f308dfdd7fd3ee717a65a4e19b31b460",
			"159973849b349fa9b37d2287af8dd528",
		},
		"sparkle/without_namespaces.xml": {
			"ee2d28f74e7d557bd7259c0f24a261658a9f27a710308a5c539ab761dae487c1",
			"d4c80271cfff4ab0afc15f7699c2e376",
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
