package repository

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"errors"
	"time"

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

func (r *FacultyRepository) loadAssociations(facs []model.Faculty) error {
	if len(facs) == 0 {
		return nil
	}

	facIDs := make([]uint, len(facs))
	for i, f := range facs {
		facIDs[i] = f.ID
	}

	var rows []dto.FacultyFlatRow

	err := r.db.Raw(`
		SELECT
			f.id AS faculty_id,
			f.name AS faculty_name,
			f.gender AS faculty_gender,
			f.joining_date,
			f.department_id,
			f.is_active AS faculty_active,

			s.id AS student_id,
			s.name AS student_name,
			s.email AS student_email,
			s.gender AS student_gender,
			s.is_active AS student_active,

			fe.id AS fee_id,
			fe.payment_mode,
			fe.amount AS fee_amount,
			fe.is_active AS fee_active

		FROM faculties f

		LEFT JOIN students s
			ON s.faculty_id = f.id
			AND s.deleted_at IS NULL

		LEFT JOIN fees fe
			ON fe.student_id = s.id
			AND fe.deleted_at IS NULL

		WHERE f.id IN ?
	`, facIDs).Scan(&rows).Error

	if err != nil {
		return err
	}


	// map faculty index
	facMap := make(map[uint]int)

	for i := range facs {
		facMap[facs[i].ID] = i
	}


	// prevent duplicate students
	studIndex := make(map[uint]map[uint]int)


	for _, row := range rows {

		facIdx, exists := facMap[row.FacultyID]
		if !exists {
			continue
		}


		// no student
		if row.StudentID == nil {
			continue
		}


		if studIndex[row.FacultyID] == nil {
			studIndex[row.FacultyID] = make(map[uint]int)
		}


		studIdx, exists := studIndex[row.FacultyID][*row.StudentID]


		// create student
		if !exists {

			student := model.Student{
				ID:        *row.StudentID,
				FacultyID: row.FacultyID,
			}


			if row.StudentName != nil {
				student.Name = *row.StudentName
			}

			if row.StudentEmail != nil {
				student.Email = *row.StudentEmail
			}

			if row.StudentGender != nil {
				student.Gender = *row.StudentGender
			}

			if row.StudentActive != nil {
				student.IsActive = *row.StudentActive
			}


			facs[facIdx].Students = append(
				facs[facIdx].Students,
				student,
			)


			studIdx = len(facs[facIdx].Students)-1

			studIndex[row.FacultyID][*row.StudentID] = studIdx
		}


		// add fee
		if row.FeeID != nil {

			student := &facs[facIdx].Students[studIdx]


			fee := model.Fees{
				ID:        *row.FeeID,
				StudentID: *row.StudentID,
			}


			if row.FeePaymentMode != nil {
				fee.PaymentMode = *row.FeePaymentMode
			}

			if row.FeeAmount != nil {
				fee.Amount = *row.FeeAmount
			}

			if row.FeeActive != nil {
				fee.IsActive = *row.FeeActive
			}


			student.Fees = append(student.Fees, fee)
		}
	}


	return nil
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
		faculty.CreatedAt = now
		faculty.UpdatedAt = now
		faculty.IsActive = true
	}

	if len(faculty.Students) > 0 {
		for i := range faculty.Students {
			faculty.Students[i].FacultyID = faculty.ID
			faculty.Students[i].CreatedAt = now
			faculty.Students[i].UpdatedAt = now
			faculty.Students[i].IsActive = true

			studRes, studErr := db.Exec(
				`INSERT INTO students (name, email, gender, faculty_id, created_at, updated_at, is_active)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				faculty.Students[i].Name,
				faculty.Students[i].Email,
				faculty.Students[i].Gender,
				faculty.Students[i].FacultyID,
				faculty.Students[i].CreatedAt,
				faculty.Students[i].UpdatedAt,
				faculty.Students[i].IsActive,
			)
			if studErr != nil {
				return studErr
			}
			studID, studIDErr := studRes.LastInsertId()
			if studIDErr == nil {
				faculty.Students[i].ID = uint(studID)
			}
		}
	}

	return nil
}

func (r *FacultyRepository) FetchFaculty() ([]model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE deleted_at IS NULL").Scan(&facs).Error
	if err != nil {
		return nil, err
	}
	err = r.loadAssociations(facs)
	return facs, err
}

func (r *FacultyRepository) FetchFacultyPaginated(page, limit int) ([]model.Faculty, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM faculties WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var facs []model.Faculty
	err = r.db.Raw("SELECT * FROM faculties WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset).Scan(&facs).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.loadAssociations(facs)
	return facs, total, err
}

func (r *FacultyRepository) FetchFacultyById(id uint) (model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&facs).Error
	if err != nil {
		return model.Faculty{}, err
	}
	if len(facs) == 0 {
		return model.Faculty{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(facs)
	if err != nil {
		return model.Faculty{}, err
	}
	return facs[0], nil
}

func (r *FacultyRepository) FetchFacultyDeleted() ([]model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE deleted_at IS NOT NULL").Scan(&facs).Error
	if err != nil {
		return nil, err
	}
	err = r.loadAssociations(facs)
	return facs, err
}

func (r *FacultyRepository) GetActiveFaculty() (model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true).Scan(&facs).Error
	if err != nil {
		return model.Faculty{}, err
	}
	if len(facs) == 0 {
		return model.Faculty{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(facs)
	if err != nil {
		return model.Faculty{}, err
	}
	return facs[0], nil
}

func (r *FacultyRepository) GetInactiveFaculty() (model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false).Scan(&facs).Error
	if err != nil {
		return model.Faculty{}, err
	}
	if len(facs) == 0 {
		return model.Faculty{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(facs)
	if err != nil {
		return model.Faculty{}, err
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
