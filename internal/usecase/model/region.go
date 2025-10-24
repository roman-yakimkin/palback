package model

import "palback/internal/domain/model"

type RegionDetail struct {
	ID      int
	Country model.Country
	Name    string
}

func CreateRegionDetail(region model.Region, country model.Country) RegionDetail {
	return RegionDetail{
		ID:      region.ID,
		Name:    region.Name,
		Country: country,
	}
}

type RegionList struct {
	Items []RegionDetail
}

func CreateRegionList(regions []model.Region, countries []model.Country) (result RegionList) {
	countryMap := make(map[string]model.Country)
	for _, country := range countries {
		countryMap[country.ID] = country
	}

	for _, region := range regions {
		result.Items = append(result.Items, CreateRegionDetail(region, countryMap[region.CountryID]))
	}

	return
}
