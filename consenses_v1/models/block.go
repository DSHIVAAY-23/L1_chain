package models


type Block struct {
	BlockHeader Header
	MetaData        []string
}

type Header struct {
	Merkleroot      string
	Datahash        string
	Prevhash        string
	Proposeraddress string
	Timestamp       string
	Height          int
}