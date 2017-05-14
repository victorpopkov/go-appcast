package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecksum(t *testing.T) {
	c := NewChecksum(Sha256, "test")
	assert.IsType(t, Checksum{}, *c)
	assert.Equal(t, Sha256, c.Algorithm)
	assert.Equal(t, "test", c.Source)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", c.Result)
}

func TestGenerate(t *testing.T) {
	testCases := map[string][]string{
		"sourceforge_default.xml": {
			"cf45ae9ba4be292c198c30663bd6bf76e3b66260b2675d5a699f10953c251288",
			"1eed329e29aa768b242d23361adf225a654e7df74d58293a44d14862ef7ef975",
			"75b31fefbd17e918078477236035a54a",
		},
		"sourceforge_empty.xml": {
			"b6ee64001ab00dbedea8fede21abb78d011c36e38bbd5aaa1004872df170c022",
			"568863d4a2540349db3987320525303f7cdd26bba6e0cada704ce2191afc9ae5",
			"30eaf5f22d3fa94b017581681886b77a",
		},
		"sourceforge_single.xml": {
			"fb59e0dba21bb8ec56d73de0f9af56547ed1951842eb682dddcb1ce453ee5443",
			"aae4e241300ef6abaf1d855b3acc613344541207159cf85064124f0a207e37ab",
			"1c177e9949f45af03df6bb83e4eeb979",
		},
		"sparkle_attributes_as_elements.xml": {
			"898628bcbf1005995c4a1e8200f6336da11fae771fc724f8fc7a9cfde8f4e85e",
			"06a16fc0d5c7f8e18ca04dbc52138159b5438cdb929e033dae6ddebca7e710fc",
			"05d4e5b0b4d005e3512a7bc24bb94925",
		},
		"sparkle_default_asc.xml": {
			"9e319d5eb9929ea069a7db81d8b46e403f05ada0dec5a4601c552a2ab08cca27",
			"8ad0cd8d67f12ed75fdfbf74e904ef8b82084875c959bec00abd5a166c512b5d",
			"da2bc13c30e16a585c0a012bcae110d5",
		},
		"sparkle_default.xml": {
			"3401290b3e7d32d01653c10668a34d53862d81f8046d7e5988bdd8b54443c2c4",
			"583743f5e8662cb223baa5e718224fa11317b0983dbf8b3c9c8d412600b6936c",
			"279ea1e0dc339ef3d04a1b9e4fd4dd82",
		},
		"sparkle_incorrect_namespace.xml": {
			"798f122b491661373cc207753dd7571590bb910860ce57ca9f3ee1ed2f9e197c",
			"f7ced8023765dc7f37c3597da7a1f8d33b3c22cc764e329babd3df16effdd245",
			"b473e0071d84b60d516e518a111d849f",
		},
		"sparkle_invalid_version.xml": {
			"5678aee518c7aaeed32bf8b8ff836d946e3baa415bc36824cc6bf4c90a96d7f3",
			"ac8bf225fb789f8174fccf26b52cde07b884e26e89546ab3ad9433cbe38ecb20",
			"a3d2cb7053b25a811f216d486469f30a",
		},
		"sparkle_multiple_enclosure.xml": {
			"48fc8531b253c5d3ed83abfe040edeeafb327d103acbbacf12c2288769dc80b9",
			"6ba0ab0e37d4280803ff2f197aaf362a3553849fb296a64bc946eda1bdb759c7",
			"9f1c1a667efc3080f1dcf020eca97c7b",
		},
		"sparkle_no_releases.xml": {
			"65911706576dab873c2b30b2d6505581d17f8e2c763da7320cfb06bbc2d4eaca",
			"65911706576dab873c2b30b2d6505581d17f8e2c763da7320cfb06bbc2d4eaca",
			"f63b85384e4c7fff3ebc14017d2edcdd",
		},
		"sparkle_single.xml": {
			"12be6a3f8d15a049e030ea09c176321278867c05742f5c2cd87aa2c368b11713",
			"98c94ba87d4eb1d99b56652b537a26d3c68efa7efa5f497839a3832a31147a7a",
			"ce7be28ec30341d08d0b4b6f24ea5c28",
		},
		"sparkle_without_comments.xml": {
			"fecbdf715eef8e743cd720d1d7799e12d569349228d3d3357cb47fee0532fec3",
			"88ceb464f652d7bf43f351f41637facd671f8f04e9a32b4b077886d24251e472",
			"f54aa1aaf762e95f86ec768f7c2e98c3",
		},
		"sparkle_without_namespaces.xml": {
			"a3f5c793c6e72f6b2cf5a24b35d8bb26b424441b22a6186c81ddc508fe0f2ae2",
			"d4cdd55c6dbf944d03c5267f3f7be4a9f7c2f1b94929359ce7e21aeef3b0747b",
			"89d619d29be8e5b03fed41465b22591e",
		},
		"unknown.xml": {
			"a4161a72df970e6fca434e2b9e256b850f12d2934cdde057985b77ea892f35d8",
			"a4161a72df970e6fca434e2b9e256b850f12d2934cdde057985b77ea892f35d8",
			"492a7260f7f43fef03d72ccffd4d27bf",
		},
	}

	for filename, checkpoints := range testCases {
		content := string(getTestdata(filename))

		// SHA256
		c := &Checksum{Sha256, content, ""}
		assert.Equal(t, checkpoints[0], c.Generate(), fmt.Sprintf("Checksum doesn't match (Sha256): %s", filename))

		// SHA256 (Homebrew-Cask)
		c = &Checksum{Sha256HomebrewCask, content, ""}
		assert.Equal(t, checkpoints[1], c.Generate(), fmt.Sprintf("Checksum doesn't match (Sha256HomebrewCask): %s", filename))

		// MD5
		c = &Checksum{Md5, content, ""}
		assert.Equal(t, checkpoints[2], c.Generate(), fmt.Sprintf("Checksum doesn't match (Md5): %s", filename))
	}
}
