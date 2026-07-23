package services

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/repository"
)

type PermissionService struct {
	permrepo *repository.PermissionRepository
}

func NewPermissionService(permrepo *repository.PermissionRepository) *PermissionService {
	return &PermissionService{permrepo: permrepo}
}

func (s *PermissionService) GetAllPermissions() ([]dto.PermissionResponseDTO, error) {
	perms, err := s.permrepo.GetAllPermissions()
	if err != nil {
		return nil, err
	}
	return dto.ToPermissionResponseListDTO(perms), nil
}
