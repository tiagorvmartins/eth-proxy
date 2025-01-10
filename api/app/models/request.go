package models

type Request struct {
	Path    string	      `json:"Path" binding:"omitempty"`
	JsonRpc string        `json:"jsonrpc" binding:"required"`
	Method  string        `json:"method" binding:"required"`
	Params  []interface{} `json:"params"`
	Id      *int          `json:"id" binding:"required"`
}
