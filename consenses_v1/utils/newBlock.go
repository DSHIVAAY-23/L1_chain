package utils

import (
	"pop_v1/config"
	"pop_v1/models"

	"time"
)

func NewBlock(txns []models.Transaction, metadata []string) models.Block {
	// generate the merkle root
	merkleRoot := GenerateMerkleRoot(txns)
	// generate the DataHash
	datahash := GenerateDataHash(txns)
	// find the  previousHash
	Hash := PreviousHash()
	// get the address of current node
	address := config.Config("HOST") + ":" + config.Config("PORT")
	// get the current time
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	//generate the header
	header := models.Header{
		Merkleroot:      merkleRoot,
		Datahash:        datahash,
		Prevhash:        Hash,
		Proposeraddress: address,
		Timestamp:       currentTime,
		Height:          len(txns),
	}

	block := models.Block{
		BlockHeader: header,
		MetaData:    metadata,
	}
	return block
}
