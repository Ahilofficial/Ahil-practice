package services

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
)

type RoleService struct {
	rolerepo *repository.RoleRepository
}

func NewRoleService(rolerepo *repository.RoleRepository) *RoleService {
	return &RoleService{rolerepo: rolerepo}
}

func (s *RoleService) CreateRole(createDTO *dto.CreateRoleDTO) (dto.RoleResponseDTO, error) {
	role := model.Role{
		Name: createDTO.Name,
	}
	err := s.rolerepo.CreateRole(&role)
	if err != nil {
		return dto.RoleResponseDTO{}, err
	}
	return dto.ToRoleResponseDTO(&role), nil
}



func (s *RoleService) GetRoleByID(id uint) (dto.RoleResponseDTO, error) {
	role, err := s.rolerepo.GetRoleByID(id)
	if err != nil {
		return dto.RoleResponseDTO{}, err
	}
	return dto.ToRoleResponseDTO(&role), nil
}



func (s *RoleService) AssignPermissionsToRole(roleID uint, assignDTO *dto.AssignPermissionsDTO) error {
	return s.rolerepo.AssignPermissionsToRole(roleID, assignDTO.PermissionIDs, assignDTO.PermissionNames)
}

func (s *RoleService) GetRolePermissions(roleID uint) ([]dto.PermissionResponseDTO, error) {
	perms, err := s.rolerepo.GetRolePermissions(roleID)
	if err != nil {
		return nil, err
	}
	return dto.ToPermissionResponseListDTO(perms), nil
}

func (s *RoleService) RemovePermissionFromRole(roleID uint, permissionID uint) error {
	return s.rolerepo.RemovePermissionFromRole(roleID, permissionID)
}
