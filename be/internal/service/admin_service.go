package service

import (
	"log"
	"milestone3/be/internal/dto"
)

type AdminRepository interface {
	CountPayment() (count int64, err error)
	CountDonation() (count int64, err error)
	CountArticle() (count int64, err error)
	CountAuction() (count int64, err error)
}

type AdminServ struct {
	adminRepo AdminRepository
}

func NewAdminService(ar AdminRepository) *AdminServ {
	return &AdminServ{adminRepo: ar}
}

func (as *AdminServ) AdminDashboard() (resp dto.AdminDashboardResponse, err error) {
	log.Println("article")
	article, err := as.adminRepo.CountArticle()
	if err != nil {
		log.Printf("error count article %s", err)
		return dto.AdminDashboardResponse{}, err
	}

	log.Println("donation")
	donation, err := as.adminRepo.CountDonation() 
	if err != nil {
		log.Printf("error count donation %s", err)
		return dto.AdminDashboardResponse{}, err
	}

	log.Println("payment")
	payment, err := as.adminRepo.CountPayment()
	if err != nil {
		log.Printf("error count payment %s", err)
		return dto.AdminDashboardResponse{}, err
	}

	auction, err := as.adminRepo.CountAuction(); 
	if err != nil {
		log.Printf("error count payment %s", err)
		return dto.AdminDashboardResponse{}, err
	}

	respon := dto.AdminDashboardResponse{
		TotalArticle: article,
		TotalDonation: donation,
		TotalPayment: payment,
		TotalAuction: auction,
	}

	return respon, nil
}

// work in progress (WIP)
// func (as *AdminServ) AdminReport() (err error) { }