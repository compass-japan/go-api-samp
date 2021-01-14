package dto

/*
 * dtoとしてシステムの入出力で扱う構造体定義
 */

type (
	AuthHeader struct {
		AuthToken string
	}

	RegisterRequest struct {
		Date    string `json:"date" validate:"required,len=8,alphanum"`
		Weather uint32 `json:"weather" validate:"required,gte=0,lte=3"`
		comment string `json:"comment"`
	}
	RegisterResponse struct {
		message string `json:"message"`
	}

	GetWeatherRequest struct {
		LocationId string `json:"-" param:"locationId" validate:"required,alphanum"`
		Date       string `json:"-" param:"date" validate:"required,len=8,alphanum"`
	}
	GetWeatherResponse struct {
		Date    string `json:"date"`
		Weather string `json:"weather"`
		Comment string `json:"comment"`
	}

	ApiDataResponse struct {
		Date    string `json:"date"`
		weather string `json:"weather"`
	}
)
