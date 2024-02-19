package services

import "github.com/tiagokriok/oka/internal/repositories"

type LinkService struct {
	repo *repositories.LinkRepository
}

func NewLinkService(repo *repositories.LinkRepository) *LinkService {
	return &LinkService{
		repo,
	}
}

func (linkSvc *LinkService) Create(link *repositories.Link) (*repositories.Link, error) {
	link, err := linkSvc.repo.CreateLink(link)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (linkSvc *LinkService) GetLinkByKey(key string) (*repositories.Link, error) {
	link, err := linkSvc.repo.GetLinkByKey(key)
	if err != nil {
		return nil, err
	}
	return link, err
}
