package repository

import (
	"backend_institutions/internal/model"
	"errors"
	"time"
	"backend_institutions/internal/dto"

	"gorm.io/gorm"
)

type FacultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) *FacultyRepository {
	return &FacultyRepository{
		db: db,
	}
}



func (r *FacultyRepository) fetchWithRelations(baseQuery string, args ...interface{}) ([]model.Faculty, error) {
	var rows []dto.FacultyFlatRow
	query := `
	SELECT 
		f.id AS fac_id, f.name AS fac_name, f.gender AS fac_gender, f.joining_date AS fac_joining_date, f.department_id, f.is_active AS fac_active,
		s.id AS stud_id, s.name AS stud_name, s.email AS stud_email, s.gender AS stud_gender, s.is_active AS stud_active,
		fe.id AS fee_id, fe.payment_mode AS fee_payment_mode, fe.amount AS fee_amount, fe.is_active AS fee_active
	FROM (` + baseQuery + `) f
	LEFT JOIN students s ON s.faculty_id = f.id AND s.deleted_at IS NULL
	LEFT JOIN fees fe ON fe.student_id = s.id AND fe.deleted_at IS NULL
	ORDER BY f.id, s.id, fe.id`

	err := r.db.Raw(query, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	facMap := make(map[uint]*model.Faculty)
	studMap := make(map[uint]*model.Student)

	var orderedIDs []uint
	for _, row := range rows {
		fac, exists := facMap[row.FacID]
		if !exists {
			fac = &model.Faculty{
				ID:           row.FacID,
				Name:         row.FacName,
				Gender:       row.FacGender,
				JoiningDate:  row.FacJoiningDate,
				DepartmentID: row.DepartmentID,
				IsActive:     row.FacActive,
				Students:     []model.Student{},
			}
			facMap[row.FacID] = fac
			orderedIDs = append(orderedIDs, row.FacID)
		}

		if row.StudID != nil {
			stud, exists := studMap[*row.StudID]
			if !exists {
				stud = &model.Student{
					ID:        *row.StudID,
					Name:      *row.StudName,
					Email:     *row.StudEmail,
					Gender:    *row.StudGender,
					FacultyID: row.FacID,
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

	result := make([]model.Faculty, len(orderedIDs))
	for i, id := range orderedIDs {
		result[i] = *facMap[id]
	}
	return result, nil
}

func (r *FacultyRepository) CreateFaculty(faculty *model.Faculty) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	query := `
		INSERT INTO faculties (name, gender, joining_date, department_id, created_at, updated_at, is_active)
		SELECT ?, ?, ?, id, ?, ?, ? FROM departments
		WHERE id = ? AND deleted_at IS NULL AND is_active = true
		  AND NOT EXISTS (
			  SELECT 1 FROM faculties 
			  WHERE name = ? AND department_id = ? AND deleted_at IS NULL
		  )`

	res, err := db.Exec(
		query,
		faculty.Name,
		faculty.Gender,
		faculty.JoiningDate,
		now,
		now,
		true,
		faculty.DepartmentID,
		faculty.Name,
		faculty.DepartmentID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("faculty name already exists in this department, or parent department is inactive/invalid")
	}
	id, err := res.LastInsertId()
	if err == nil {
		faculty.ID = uint(id)
	}
	return nil
}

func (r *FacultyRepository) FetchFaculty() ([]model.Faculty, error) {
	return r.fetchWithRelations("SELECT * FROM faculties WHERE deleted_at IS NULL")
}

func (r *FacultyRepository) FetchFacultyPaginated(page, limit int) ([]model.Faculty, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM faculties WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	facs, err := r.fetchWithRelations("SELECT * FROM faculties WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset)
	return facs, total, err
}

func (r *FacultyRepository) FetchFacultyById(id uint) (model.Faculty, error) {
	facs, err := r.fetchWithRelations("SELECT * FROM faculties WHERE id = ? AND deleted_at IS NULL LIMIT 1", id)
	if err != nil {
		return model.Faculty{}, err
	}
	if len(facs) == 0 {
		return model.Faculty{}, gorm.ErrRecordNotFound
	}
	return facs[0], nil
}

func (r *FacultyRepository) FetchFacultyDeleted() ([]model.Faculty, error) {
	return r.fetchWithRelations("SELECT * FROM faculties WHERE deleted_at IS NOT NULL")
}

func (r *FacultyRepository) GetActiveFacutly() (model.Faculty, error) {
	facs, err := r.fetchWithRelations("SELECT * FROM faculties WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true)
	if err != nil {
		return model.Faculty{}, err
	}
	if len(facs) == 0 {
		return model.Faculty{}, gorm.ErrRecordNotFound
	}
	return facs[0], nil
}

func (r *FacultyRepository) GetInactiveFaculty() (model.Faculty, error) {
	facs, err := r.fetchWithRelations("SELECT * FROM faculties WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false)
	if err != nil {
		return model.Faculty{}, err
	}
	if len(facs) == 0 {
		return model.Faculty{}, gorm.ErrRecordNotFound
	}
	return facs[0], nil
}

func (r *FacultyRepository) DeleteFaculty(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE faculties SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
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

func (r *FacultyRepository) UpdateFacultyById(faculty *model.Faculty) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE faculties SET name = ?, gender = ?, updated_at = ? WHERE id = ?",
		faculty.Name, faculty.Gender, time.Now(), faculty.ID,
	)
	return err
}
