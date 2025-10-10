package dto

import "palback/internal/domain/model"

type CountryPostRequest struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	HasRegions bool   `json:"has_regions"`
}

type CountryPutRequest struct {
	Name       string `json:"name"`
	HasRegions bool   `json:"has_regions"`
}

type CountryResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	HasRegions bool   `json:"has_regions"`
}

func CreateCountryResponse(src model.Country) CountryResponse {
	return CountryResponse{
		ID:         src.ID,
		Name:       src.Name,
		HasRegions: src.HasRegions,
	}
}

type CountryResponseList struct {
	Items []CountryResponse `json:"items"`
}

func CreateCountryResponseList(src []model.Country) CountryResponseList {
	result := CountryResponseList{
		Items: make([]CountryResponse, 0, len(src)),
	}

	for _, country := range src {
		result.Items = append(result.Items, CreateCountryResponse(country))
	}

	return result
}
