package models

type Request struct {
	JsonRpc string        `json:"jsonrpc" binding:"required"`
	Method  string        `json:"method" binding:"required"`
	Params  []interface{} `json:"params" binding:"required"`
	Id      *int          `json:"id" binding:"required"`
}
