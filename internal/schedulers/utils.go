package schedulers

var (
	hashSize = map[string]int{
		"MD5Sum": 32,
		"SHA1":   40,
		"SHA256": 64,
		"SHA512": 128,
	}
)
