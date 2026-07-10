package repository

import (
	"backend_institutions/internal/model"
	"errors"
	"time"
	"backend_institutions/internal/dto"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}



func (r *DepartmentRepository) fetchWithRelations(baseQuery string, args ...interface{}) ([]model.Department, error) {
	var rows []dto.DepartmentFlatRow
	query := `
	SELECT 
		d.id AS dept_id, d.department_name, d.institution_id, d.is_active AS dept_active,
		f.id AS fac_id, f.name AS fac_name, f.gender AS fac_gender, f.joining_date AS fac_joining_date, f.is_active AS fac_active,
		s.id AS stud_id, s.name AS stud_name, s.email AS stud_email, s.gender AS stud_gender, s.is_active AS stud_active,
		fe.id AS fee_id, fe.payment_mode AS fee_payment_mode, fe.amount AS fee_amount, fe.is_active AS fee_active
	FROM (` + baseQuery + `) d
	LEFT JOIN faculties f ON f.department_id = d.id AND f.deleted_at IS NULL
	LEFT JOIN students s ON s.faculty_id = f.id AND s.deleted_at IS NULL
	LEFT JOIN fees fe ON fe.student_id = s.id AND fe.deleted_at IS NULL
	ORDER BY d.id, f.id, s.id, fe.id`

	err := r.db.Raw(query, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	deptMap := make(map[uint]*model.Department)
	facMap := make(map[uint]*model.Faculty)
	studMap := make(map[uint]*model.Student)

	var orderedIDs []uint
	for _, row := range rows {
		dept, exists := deptMap[row.DeptID]
		if !exists {
			dept = &model.Department{
				ID:             row.DeptID,
				DepartmentName: row.DepartmentName,
				InstitutionID:  row.InstitutionID,
				IsActive:       row.DeptActive,
				Faculties:      []model.Faculty{},
			}
			deptMap[row.DeptID] = dept
			orderedIDs = append(orderedIDs, row.DeptID)
		}

		if row.FacID != nil {
			fac, exists := facMap[*row.FacID]
			if !exists {
				fac = &model.Faculty{
					ID:           *row.FacID,
					Name:         *row.FacName,
					Gender:       *row.FacGender,
					JoiningDate:  *row.FacJoiningDate,
					DepartmentID: row.DeptID,
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

	result := make([]model.Department, len(orderedIDs))
	for i, id := range orderedIDs {
		result[i] = *deptMap[id]
	}
	return result, nil
}

func (r *DepartmentRepository) CreateDepartment(department *model.Department) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	res, err := db.Exec(
		`INSERT INTO departments (department_name, institution_id, created_at, updated_at, is_active)
		SELECT ?, id, ?, ?, ? FROM institutions
		WHERE id = ? AND deleted_at IS NULL AND is_active = true
		  AND NOT EXISTS (
			  SELECT 1 FROM departments 
			  WHERE department_name = ? AND institution_id = ? AND deleted_at IS NULL
		  )`,
		department.DepartmentName, now, now, true,
		department.InstitutionID,
		department.DepartmentName, department.InstitutionID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("department name already exists in this institution, or parent institution is inactive/invalid")
	}
	id, err := res.LastInsertId()
	if err == nil {
		department.ID = uint(id)
	}
	return nil
}

func (r *DepartmentRepository) FetchDepartment() ([]model.Department, error) {
	return r.fetchWithRelations("SELECT * FROM departments WHERE deleted_at IS NULL")
}

func (r *DepartmentRepository) FetchDepartmentPaginated(page, limit int) ([]model.Department, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM departments WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	depts, err := r.fetchWithRelations("SELECT * FROM departments WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset)
	return depts, total, err
}

func (r *DepartmentRepository) FetchDepartmentById(id uint) (model.Department, error) {
	depts, err := r.fetchWithRelations("SELECT * FROM departments WHERE id = ? AND deleted_at IS NULL LIMIT 1", id)
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	return depts[0], nil
}

func (r *DepartmentRepository) FetchDepartmentDeleted() ([]model.Department, error) {
	return r.fetchWithRelations("SELECT * FROM departments WHERE deleted_at IS NOT NULL")
}

func (r *DepartmentRepository) GetActiveDepartment() (model.Department, error) {
	depts, err := r.fetchWithRelations("SELECT * FROM departments WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true)
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	return depts[0], nil
}

func (r *DepartmentRepository) GetInactiveDepartment() (model.Department, error) {
	depts, err := r.fetchWithRelations("SELECT * FROM departments WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false)
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	return depts[0], nil
}

func (r *DepartmentRepository) DeleteDepartment(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE departments SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
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

func (r *DepartmentRepository) UpdateDepartmentById(department *model.Department) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE departments SET department_name = ?, updated_at = ? WHERE id = ?",
		department.DepartmentName, time.Now(), department.ID,
	)
	return err
}
