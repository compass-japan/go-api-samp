package dto

/*
 * dtoとしてシステムの入出力で扱う構造体定義
 */

type (
	RegisterRequest struct {
		LocationId int    `json:"location_id" validate:"required"`
		Date       string `json:"date" validate:"required,len=8,alphanum"`
		Weather    int    `json:"weather" validate:"required,min=0,max=3"`
		Comment    string `json:"comment" validate:"lte=255"`
	}
	RegisterResponse struct {
		Message string `json:"message"`
	}

	GetWeatherRequest struct {
		LocationId int    `json:"-" param:"locationId" validate:"required"`
		Date       string `json:"-" param:"date" validate:"required,len=8,alphanum"`
	}
	GetWeatherResponse struct {
		Location string `json:"location"`
		Date     string `json:"date"`
		Weather  string `json:"weather"`
		Comment  string `json:"comment"`
	}

	ExApiResponse struct {
		ConsolidatedWeather []ConsolidatedWeather `json:"consolidated_weather"`
		Time                string                `json:"-"` //`json:"time"`
		SunRize             string                `json:"-"` //`json:"sun_rize"`
		SunSet              string                `json:"-"` //`json:"sun_set"`
		TimezoneName        string                `json:"-"` //`json:"timezone_name"`
		Parent              string                `json:"-"` //`json:"parent"` #todo object
		Sources             []string              `json:"-"` //`json:"sources"` #todo object
		Title               string                `json:"title"`
		LocationType        string                `json:"-"` //`json:"location_type"`
		Woeid               int                   `json:"-"` //`json:"woeid"`
		LattLong            string                `json:"-"` //`json:"latt_long"`
		Timezone            string                `json:"timezone"`
	}

	ConsolidatedWeather struct {
		Id                   int     `json:"-"` //`json:"id"`
		WeatherStateName     string  `json:"weather_state_name"`
		WeatherStateAbbr     string  `json:"-"` //`json:"weather_state_abbr":`
		WindDirectionCompass string  `json:"-"` //`json:wind_direction_compass"`
		Created              string  `json:"-"` //`json:"created"`
		ApplicableDate       string  `json:"applicable_date"`
		MinTemp              float64 `json:"-"` //`json:"min_temp"`
		MaxTemp              float64 `json:"-"` //`json:"max_temp"`
		TheTemp              float64 `json:"-"` //`json:"the_temp"`
		WindSpeed            float64 `json:"wind_speed"`
		WindDirection        float64 `json:"-"` //`json:"wind_direction"`
		AirPressure          float64 `json:"air_pressure"`
		Humidity             int     `json:"humidity"`
		Visibility           float64 `json:"-"` //`json:"visibility"`
		Predictability       int     `json:"-"` //`json:"predictability":`
	}
)
