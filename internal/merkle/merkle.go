package merkle

import (
	"crypto/sha256"
)

func MerkleRoot(hashes [][32]byte) [32]byte {
	merkle := hashes

	if len(merkle)%2 != 0 {
		merkle = append(merkle, merkle[len(merkle)-1])
	}

	for len(merkle) > 1 {
		var newMerkle [][32]byte
		for i := 0; i < len(merkle); i += 2 {
			mergedHashes := append(merkle[i][:], merkle[i+1][:]...)
			hash := sha256.Sum256(mergedHashes)
			newMerkle = append(newMerkle[:], hash)
		}

		merkle = newMerkle
	}

	return merkle[0]
}
