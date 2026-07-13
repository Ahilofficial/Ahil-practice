package repository

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) loadAssociations(depts []model.Department) error {
	if len(depts) == 0 {
		return nil
	}

	deptIDs := make([]uint, len(depts))
	for i, dept := range depts {
		deptIDs[i] = dept.ID
	}

	var rows []dto.DepartmentFlatRow

	err := r.db.Raw(`
		SELECT
			d.id AS dept_id,
			d.department_name AS department_name,
			d.institution_id,
			d.is_active AS dept_active,

			f.id AS fac_id,
			f.name AS fac_name,
			f.gender AS fac_gender,
			f.joining_date AS fac_joining_date,
			f.is_active AS fac_active,

			s.id AS stud_id,
			s.name AS stud_name,
			s.email AS stud_email,
			s.gender AS stud_gender,
			s.is_active AS stud_active,

			fe.id AS fee_id,
			fe.payment_mode AS fee_payment_mode,
			fe.amount AS fee_amount,
			fe.is_active AS fee_active

		FROM departments d

		LEFT JOIN faculties f
			ON f.department_id = d.id
			AND f.deleted_at IS NULL

		LEFT JOIN students s
			ON s.faculty_id = f.id
			AND s.deleted_at IS NULL

		LEFT JOIN fees fe
			ON fe.student_id = s.id
			AND fe.deleted_at IS NULL

		WHERE d.id IN ?

	`, deptIDs).Scan(&rows).Error

	if err != nil {
		return err
	}


	deptIndex := make(map[uint]int)

	for i := range depts {
		deptIndex[depts[i].ID] = i
	}


	facIndex := make(map[uint]map[uint]int)
	studIndex := make(map[uint]map[uint]int)


	for _, row := range rows {

		deptPos, exists := deptIndex[row.DeptID]
		if !exists {
			continue
		}


		// No faculty
		if row.FacID == nil {
			continue
		}


		// Faculty map initialization
		if facIndex[row.DeptID] == nil {
			facIndex[row.DeptID] = make(map[uint]int)
		}


		facPos, exists := facIndex[row.DeptID][*row.FacID]


		// Create Faculty
		if !exists {

			fac := model.Faculty{
				ID:           *row.FacID,
				DepartmentID: row.DeptID,
			}

			if row.FacName != nil {
				fac.Name = *row.FacName
			}

			if row.FacGender != nil {
				fac.Gender = *row.FacGender
			}

			if row.FacJoiningDate != nil {
				fac.JoiningDate = *row.FacJoiningDate
			}

			if row.FacActive != nil {
				fac.IsActive = *row.FacActive
			}


			depts[deptPos].Faculties =
				append(depts[deptPos].Faculties, fac)


			facPos = len(depts[deptPos].Faculties) - 1

			facIndex[row.DeptID][*row.FacID] = facPos
		}



		// No student
		if row.StudID == nil {
			continue
		}


		faculty := &depts[deptPos].Faculties[facPos]


		if studIndex[*row.FacID] == nil {
			studIndex[*row.FacID] = make(map[uint]int)
		}


		studPos, exists := studIndex[*row.FacID][*row.StudID]


		// Create Student
		if !exists {

			student := model.Student{
				ID:        *row.StudID,
				FacultyID: *row.FacID,
			}


			if row.StudName != nil {
				student.Name = *row.StudName
			}

			if row.StudEmail != nil {
				student.Email = *row.StudEmail
			}

			if row.StudGender != nil {
				student.Gender = *row.StudGender
			}

			if row.StudActive != nil {
				student.IsActive = *row.StudActive
			}


			faculty.Students =
				append(faculty.Students, student)


			studPos = len(faculty.Students) - 1

			studIndex[*row.FacID][*row.StudID] = studPos
		}



		// Add Fee
		if row.FeeID != nil {

			student := &faculty.Students[studPos]

			fee := model.Fees{
				ID:        *row.FeeID,
				StudentID: *row.StudID,
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
	if err != nil {
		return err
	}

	department.ID = uint(id)
	department.CreatedAt = now
	department.UpdatedAt = now
	department.IsActive = true

	if len(department.Faculties) > 0 {
		for i := range department.Faculties {
			department.Faculties[i].DepartmentID = department.ID
			department.Faculties[i].CreatedAt = now
			department.Faculties[i].UpdatedAt = now
			department.Faculties[i].IsActive = true

			facRes, facErr := db.Exec(
				`INSERT INTO faculties (name, gender, joining_date, department_id, created_at, updated_at, is_active)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				department.Faculties[i].Name,
				department.Faculties[i].Gender,
				department.Faculties[i].JoiningDate,
				department.Faculties[i].DepartmentID,
				department.Faculties[i].CreatedAt,
				department.Faculties[i].UpdatedAt,
				department.Faculties[i].IsActive,
			)
			if facErr != nil {
				return facErr
			}
			facID, facIDErr := facRes.LastInsertId()
			if facIDErr != nil {
				return facIDErr
			}
			department.Faculties[i].ID = uint(facID)
		}
	}

	return nil
}

func (r *DepartmentRepository) FetchDepartment() ([]model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE deleted_at IS NULL").Scan(&depts).Error
	if err != nil {
		return nil, err
	}
	err = r.loadAssociations(depts)
	return depts, err
}

func (r *DepartmentRepository) FetchDepartmentPaginated(page, limit int) ([]model.Department, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM departments WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var depts []model.Department
	err = r.db.Raw("SELECT * FROM departments WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset).Scan(&depts).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.loadAssociations(depts)
	return depts, total, err
}

func (r *DepartmentRepository) FetchDepartmentById(id uint) (model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&depts).Error
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(depts)
	if err != nil {
		return model.Department{}, err
	}
	return depts[0], nil
}

func (r *DepartmentRepository) FetchDepartmentDeleted() ([]model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE deleted_at IS NOT NULL").Scan(&depts).Error
	if err != nil {
		return nil, err
	}
	err = r.loadAssociations(depts)
	return depts, err
}

func (r *DepartmentRepository) GetActiveDepartment() (model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true).Scan(&depts).Error
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(depts)
	if err != nil {
		return model.Department{}, err
	}
	return depts[0], nil
}

func (r *DepartmentRepository) GetInactiveDepartment() (model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false).Scan(&depts).Error
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(depts)
	if err != nil {
		return model.Department{}, err
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
