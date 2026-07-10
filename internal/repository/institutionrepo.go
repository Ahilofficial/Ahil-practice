package repository

import (
	"backend_institutions/internal/model"
	"errors"
	"time"
	"backend_institutions/internal/dto"
	"gorm.io/gorm"
)

type InstitutionRepository struct {
	db *gorm.DB
}

func NewInstitutionRepository(db *gorm.DB) *InstitutionRepository {
	return &InstitutionRepository{db: db}
}



func (r *InstitutionRepository) fetchWithRelations(baseQuery string, args ...interface{}) ([]model.Institutions, error) {
	var rows []dto.InstitutionFlatRow
	query := `
	SELECT 
		i.id AS inst_id, i.name AS inst_name, i.institution_code, i.state AS inst_state, i.is_active AS inst_active,
		d.id AS dept_id, d.department_name, d.is_active AS dept_active,
		f.id AS fac_id, f.name AS fac_name, f.gender AS fac_gender, f.joining_date AS fac_joining_date, f.is_active AS fac_active,
		s.id AS stud_id, s.name AS stud_name, s.email AS stud_email, s.gender AS stud_gender, s.is_active AS stud_active,
		fe.id AS fee_id, fe.payment_mode AS fee_payment_mode, fe.amount AS fee_amount, fe.is_active AS fee_active
	FROM (` + baseQuery + `) i
	LEFT JOIN departments d ON d.institution_id = i.id AND d.deleted_at IS NULL
	LEFT JOIN faculties f ON f.department_id = d.id AND f.deleted_at IS NULL
	LEFT JOIN students s ON s.faculty_id = f.id AND s.deleted_at IS NULL
	LEFT JOIN fees fe ON fe.student_id = s.id AND fe.deleted_at IS NULL
	ORDER BY i.id, d.id, f.id, s.id, fe.id`

	err := r.db.Raw(query, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	instMap := make(map[uint]*model.Institutions)
	deptMap := make(map[uint]*model.Department)
	facMap := make(map[uint]*model.Faculty)
	studMap := make(map[uint]*model.Student)

	var orderedIDs []uint
	for _, row := range rows {
		inst, exists := instMap[row.InstID]
		if !exists {
			inst = &model.Institutions{
				ID:              row.InstID,
				Name:            row.InstName,
				InstitutionCode: row.InstitutionCode,
				State:           row.InstState,
				IsActive:        row.InstActive,
				Departments:     []model.Department{},
			}
			instMap[row.InstID] = inst
			orderedIDs = append(orderedIDs, row.InstID)
		}

		if row.DeptID != nil {
			dept, exists := deptMap[*row.DeptID]
			if !exists {
				dept = &model.Department{
					ID:             *row.DeptID,
					DepartmentName: *row.DepartmentName,
					InstitutionID:  row.InstID,
					IsActive:       *row.DeptActive,
					Faculties:      []model.Faculty{},
				}
				deptMap[*row.DeptID] = dept
				inst.Departments = append(inst.Departments, *dept)
				dept = &inst.Departments[len(inst.Departments)-1]
				deptMap[*row.DeptID] = dept
			}

			if row.FacID != nil {
				fac, exists := facMap[*row.FacID]
				if !exists {
					fac = &model.Faculty{
						ID:           *row.FacID,
						Name:         *row.FacName,
						Gender:       *row.FacGender,
						JoiningDate:  *row.FacJoiningDate,
						DepartmentID: *row.DeptID,
						IsActive:     *row.FacActive,
						Students:     []model.Student{},
					}
					facMap[*row.FacID] = fac
					dept.Faculties = append(dept.Faculties, *fac)
					fac = &dept.Faculties[len(dept.Faculties)-1]
					facMap[*row.FacID] = fac
				}

				if row.StudID != nil {
					stud, exists := studMap[*row.StudID]
					if !exists {
						stud = &model.Student{
							ID:        *row.StudID,
							Name:      *row.StudName,
							Email:     *row.StudEmail,
							Gender:    *row.StudGender,
							FacultyID: *row.FacID,
							IsActive:  *row.StudActive,
							Fees:      []model.Fees{},
						}
						studMap[*row.StudID] = stud
						fac.Students = append(fac.Students, *stud)
						stud = &fac.Students[len(fac.Students)-1]
						studMap[*row.StudID] = stud
					}

					if row.FeeID != nil {
						feeFound := false
						for _, ef := range stud.Fees {
							if ef.ID == *row.FeeID {
								feeFound = true
								break
							}
						}
						if !feeFound {
							fee := model.Fees{
								ID:          *row.FeeID,
								PaymentMode: *row.FeePaymentMode,
								Amount:      *row.FeeAmount,
								StudentID:   *row.StudID,
								IsActive:    *row.FeeActive,
							}
							stud.Fees = append(stud.Fees, fee)
						}
					}
				}
			}
		}
	}

	result := make([]model.Institutions, len(orderedIDs))
	for i, id := range orderedIDs {
		result[i] = *instMap[id]
	}
	return result, nil
}

func (r *InstitutionRepository) CreateInstitution(institute *model.Institutions) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	res, err := db.Exec(
		`INSERT INTO institutions (name, institution_code, state, created_at, updated_at, is_active)
		SELECT ?, ?, ?, ?, ?, ? FROM DUAL
		WHERE NOT EXISTS (
			SELECT 1 FROM institutions 
			WHERE (name = ? OR institution_code = ?) AND deleted_at IS NULL
		)`,
		institute.Name, institute.InstitutionCode, institute.State, now, now, true,
		institute.Name, institute.InstitutionCode,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("institution name or code already exists")
	}
	id, err := res.LastInsertId()
	if err == nil {
		institute.ID = uint(id)
	}
	return nil
}

func (r *InstitutionRepository) FetchInstitution() ([]model.Institutions, error) {
	return r.fetchWithRelations("SELECT * FROM institutions WHERE deleted_at IS NULL")
}

func (r *InstitutionRepository) FetchInstitutionPaginated(page, limit int) ([]model.Institutions, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM institutions WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	insts, err := r.fetchWithRelations("SELECT * FROM institutions WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset)
	return insts, total, err
}

func (r *InstitutionRepository) FetchInstitutionById(id uint) (model.Institutions, error) {
	insts, err := r.fetchWithRelations("SELECT * FROM institutions WHERE id = ? AND deleted_at IS NULL LIMIT 1", id)
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	return insts[0], nil
}

func (r *InstitutionRepository) GetActiveInstitute() (model.Institutions, error) {
	insts, err := r.fetchWithRelations("SELECT * FROM institutions WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true)
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	return insts[0], nil
}

func (r *InstitutionRepository) GetInactiveInstitute() (model.Institutions, error) {
	insts, err := r.fetchWithRelations("SELECT * FROM institutions WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false)
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	return insts[0], nil
}

func (r *InstitutionRepository) FetchInstitutionDeleted() ([]model.Institutions, error) {
	return r.fetchWithRelations("SELECT * FROM institutions WHERE deleted_at IS NOT NULL")
}

func (r *InstitutionRepository) DeleteInstitution(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE institutions SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
		false, time.Now(), id, true,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("record not found or already deleted")
	}
	return nil
}

func (r *InstitutionRepository) UpdateInstitution(institute *model.Institutions) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE institutions SET name = ?, institution_code = ?, state = ?, updated_at = ? WHERE id = ?",
		institute.Name, institute.InstitutionCode, institute.State, time.Now(), institute.ID,
	)
	return err
}
