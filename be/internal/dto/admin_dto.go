package dto

type TotalArticle struct {
	Count int64	
}

type TotalDonation struct {
	Count int64
}

type TotalPayment struct {
	Count int64
}

type TotalAuction struct {
	Count int64 
}

type AdminDashboardResponse struct {
	TotalArticle int64 `json:"total_article"`
	TotalDonation int64 `json:"total_donation"`
	TotalPayment int64 `json:"total_payment"`
	TotalAuction int64 `json:"total_auction"`
}

// type AdminReportResponse struct { 

// }