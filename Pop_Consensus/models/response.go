package models

type Response struct {
	Block Block
	Vote bool
	Ip string
}

var VoteMap =make(map[string][]Response)
//string->[true,false]