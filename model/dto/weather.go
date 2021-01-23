package dto

/*
 * dtoとしてシステムの入出力で扱う構造体定義
 */

type (
	AuthHeader struct {
		AuthToken string
	}

	RegisterRequest struct {
		LocationId uint32 `json:"location_id" validate:"required"`
		Date       string `json:"date" validate:"required,len=8,alphanum"`
		Weather    uint32 `json:"weather" validate:"required,min=0,max=3"`
		Comment    string `json:"comment" validate:"lte=255"`
	}
	RegisterResponse struct {
		Message string `json:"message"`
	}

	GetWeatherRequest struct {
		LocationId uint32 `json:"-" param:"locationId" validate:"required"`
		Date       string `json:"-" param:"date" validate:"required,len=8,alphanum"`
	}
	GetWeatherResponse struct {
		LocationId uint32 `json:"-" param:"locationId" validate:"required"`
		Date       string `json:"date"`
		Weather    string `json:"weather"`
		Comment    string `json:"comment"`
	}

	ApiDataResponse struct {
		Date    string `json:"date"`
		Weather string `json:"weather"`
	}
)
