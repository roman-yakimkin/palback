package dto

import "palback/internal/domain/model"

type PlaceTypeResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

func CreatePlaceTypeResponse(src model.PlaceType) PlaceTypeResponse {
	return PlaceTypeResponse{
		ID:     src.ID,
		Name:   src.Name,
		Weight: src.Weight,
	}
}

type PlaceTypeResponseList struct {
	Items []PlaceTypeResponse `json:"items"`
}

func CreatePlaceTypeResponseList(src []model.PlaceType) PlaceTypeResponseList {
	result := PlaceTypeResponseList{
		Items: make([]PlaceTypeResponse, 0, len(src)),
	}

	for _, placeType := range src {
		result.Items = append(result.Items, CreatePlaceTypeResponse(placeType))
	}

	return result
}
