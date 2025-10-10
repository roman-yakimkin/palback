package dto

import appModel "palback/internal/app/model"

type RegionPostRequest struct {
	CountryID string `json:"country_id"`
	Name      string `json:"name"`
}

type RegionPutRequest struct {
	CountryID string `json:"country_id"`
	Name      string `json:"name"`
}

type RegionResponse struct {
	ID      int             `json:"id"`
	Name    string          `json:"name"`
	Country CountryResponse `json:"country"`
}

func CreateRegionResponse(src appModel.RegionDetail) RegionResponse {
	return RegionResponse{
		ID:      src.ID,
		Name:    src.Name,
		Country: CreateCountryResponse(src.Country),
	}
}

type RegionResponseList struct {
	Items []RegionResponse `json:"items"`
}

func CreateRegionResponseList(src appModel.RegionList) RegionResponseList {
	result := RegionResponseList{
		Items: make([]RegionResponse, 0, len(src.Items)),
	}

	for _, item := range src.Items {
		result.Items = append(result.Items, CreateRegionResponse(item))
	}

	return result
}
