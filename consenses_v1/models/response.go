package models

type Response struct {
	Block Block
	Vote bool
	Id string
}

var VoteMap =make(map[string][]Response)
//string->[true,false]