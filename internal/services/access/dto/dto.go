package dto

import (
	entity "github.com/av-ugolkov/lingua-evo/internal/services/access"
)

type AccessRs struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func AccessesToDto(accesses []entity.Access) []AccessRs {
	accessesRs := make([]AccessRs, 0, len(accesses))
	for _, access := range accesses {
		accessesRs = append(accessesRs, AccessRs{
			ID:   access.ID,
			Type: access.Type,
			Name: access.Name,
		})
	}

	return accessesRs
}
