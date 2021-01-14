package dto

/*
 * dtoとしてシステムのエラー出力で扱う構造体定義
 */

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}
)
