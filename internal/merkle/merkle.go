package merkle

import (
	"crypto/sha256"
	"fmt"
)

func MerkleRoot(hashes []string) string {
	merkle := hashes

	if len(hashes) == 0 {
		return "0"
	}

	if len(merkle)%2 != 0 {
		merkle = append(merkle, merkle[len(merkle)-1])
	}

	for len(merkle) > 1 {
		var newMerkle []string
		fmt.Println(merkle)
		for i := 0; i < len(merkle); i += 2 {
			hash := sha256.Sum256([]byte(merkle[i] + merkle[i+1]))
			newMerkle = append(newMerkle, fmt.Sprintf("%x", hash))
		}

		merkle = newMerkle
	}

	return merkle[0]
}
