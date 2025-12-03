package service

import "errors"

var (
	// User Errors
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidUser          = errors.New("invalid user data")
	ErrUserNotFoundID       = errors.New("user ID not found")
	ErrUserNotFoundName     = errors.New("user name not found")
	ErrUserNotFoundEmail    = errors.New("user email not found")
	ErrUserNotFoundPassword = errors.New("user password not found")
	// Payment Errors
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrInvalidPayment        = errors.New("invalid payment data")
	ErrPaymentNotFoundID     = errors.New("payment ID not found")
	ErrPaymentNotFoundAmount = errors.New("payment amount not found")
	ErrPaymentNotFoundMethod = errors.New("payment method not found")
	ErrPaymentNotFoundStatus = errors.New("payment status not found")
	// Auction Errors
	ErrAuctionNotFound      = errors.New("auction not found")
	ErrInvalidAuction       = errors.New("invalid auction data")
	ErrAuctionNotFoundID    = errors.New("auction ID not found")
	ErrSessionNotFoundID    = errors.New("auction session ID not found")
	ErrAuctionNotFoundTitle = errors.New("auction title not found")
	ErrAuctionNotFoundDesc  = errors.New("auction description not found")
	ErrAuctionNotFoundStart = errors.New("auction start time not found")
	ErrAuctionNotFoundEnd   = errors.New("auction end time not found")
	ErrInvalidDate          = errors.New("end time must be after start time")
	ErrActiveSession        = errors.New("cannot modify an active auction session")
	// Donation Errors
	ErrDonationNotFound          = errors.New("donation not found")
	ErrInvalidDonation           = errors.New("invalid donation data")
	ErrDonationNotFoundID        = errors.New("donation ID not found")
	ErrDonationNotFoundAmount    = errors.New("donation amount not found")
	ErrDonationNotFoundDonorName = errors.New("donor name not found")
	// Article Errors
	ErrArticleNotFound = errors.New("article not found")
	ErrInvalidArticle  = errors.New("invalid article data")
	// Bidding Errors
	ErrBiddingNotFound       = errors.New("bidding not found")
	ErrInvalidBidding        = errors.New("invalid bidding data")
	ErrBiddingNotFoundID     = errors.New("bidding ID not found")
	ErrBiddingNotFoundAmount = errors.New("bidding amount not found")
	ErrBidTooLow             = errors.New("bid too low")
	// Final Donation Errors
	ErrFinalDonationNotFound   = errors.New("final donation not found")
	ErrFinalDonationNotFoundID = errors.New("final donation ID not found")
	// image Errors
	ErrImageNotFound   = errors.New("image not found")
	ErrSignedURLFailed = errors.New("signed URL generation failed")

	// Authorization / Generic Errors
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden    = errors.New("forbidden access")
)
