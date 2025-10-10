package dto

import "palback/internal/domain/model"

type CityTypeResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Weight    int    `json:"weight"`
}

func CreateCityTypeResponse(src model.CityType) CityTypeResponse {
	return CityTypeResponse{
		ID:        src.ID,
		Name:      src.Name,
		ShortName: src.ShortName,
		Weight:    src.Weight,
	}
}

type CityTypeResponseList struct {
	Items []CityTypeResponse `json:"items"`
}

func CreateCityTypeResponseList(src []model.CityType) CityTypeResponseList {
	result := CityTypeResponseList{
		Items: make([]CityTypeResponse, 0, len(src)),
	}

	for _, cityType := range src {
		result.Items = append(result.Items, CreateCityTypeResponse(cityType))
	}

	return result
}
