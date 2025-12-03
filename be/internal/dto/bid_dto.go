package dto

import (
	"milestone3/be/internal/entity"
)

type BidDTO struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

func BidRequest(d BidDTO) entity.Bid {
	return entity.Bid{
		Amount: d.Amount,
	}
}

func BidResponse(m entity.Bid) BidDTO {
	return BidDTO{
		Amount: m.Amount,
	}
}

func BidResponses(ms []entity.Bid) []BidDTO {
	res := make([]BidDTO, 0, len(ms))
	for _, m := range ms {
		res = append(res, BidResponse(m))
	}
	return res
}
