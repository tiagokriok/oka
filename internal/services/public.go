package services

import "github.com/tiagokriok/oka/internal/repositories"

type PublicService struct {
	repo *repositories.LinkRepository
}

func NewPublicService(repo *repositories.LinkRepository) *PublicService {
	return &PublicService{
		repo,
	}
}

func (publicSvc *PublicService) GetLinkByKey(key string) (*repositories.Link, error) {
	link, err := publicSvc.repo.GetLinkByKey(key)
	if err != nil {
		return nil, err
	}
	return link, err
}
