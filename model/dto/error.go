package dto

/*
 * dtoとしてシステムのエラーレスポンスの構造体定義
 */

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}
)
