package app

import (


	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
)

type TopupUsecase interface {
	CreateTopup(topup *domain.TopUp) error
	GetTopupByID(topupID int) (*domain.TopUp, error)
	UpdateTopup(topup *domain.TopUp) error
	DeleteTopup(topupID int) error
	HistoryTransaction(topupID int) ([]*domain.TopUp, error)
	GetLastTopupAmount(walletID int) (float64, error)
}

type topupUsecase struct {
	topupRepository repository.TopupRepository
}

func NewTopupUsecase(topupRepository repository.TopupRepository) TopupUsecase {
	return &topupUsecase{
		topupRepository: topupRepository,
	}
}

func (u *topupUsecase) CreateTopup(topup *domain.TopUp) error {
	return u.topupRepository.Create(topup)
}

func (u *topupUsecase) GetTopupByID(topupID int) (*domain.TopUp, error) {
	return u.topupRepository.FindOne(topupID)
}

func (u *topupUsecase) UpdateTopup(topup *domain.TopUp) error {
	return u.topupRepository.Update(topup)
}

func (u *topupUsecase) DeleteTopup(topupID int) error {
	return u.topupRepository.Delete(topupID)
}

func (u *topupUsecase) GetLastTopupAmount(walletID int) (float64, error) {
	return u.topupRepository.GetLastTopupAmount(walletID)
}


func (u *topupUsecase) HistoryTransaction(topupID int) ([]*domain.TopUp, error) {
	return u.topupRepository.HistoryTopup(topupID)
}