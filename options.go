package tdam

type ContractType string

const (
	CALL ContractType = "CALL"
	PUT  ContractType = "PUT"
	ALL  ContractType = "ALL"
)
