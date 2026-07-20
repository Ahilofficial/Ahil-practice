package repository

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type InstitutionRepository struct {
	db *gorm.DB
}

func NewInstitutionRepository(db *gorm.DB) *InstitutionRepository {
	return &InstitutionRepository{
		db: db,
	}
}

func (r *InstitutionRepository) loadAssociations(insts []model.Institutions) error {
	if len(insts) == 0 {
		return nil
	}

	instIDs := make([]uint, len(insts))
	for i, inst := range insts {
		instIDs[i] = inst.ID
	}

	var rows []dto.InstitutionFlatRow

	err := r.db.Raw(`
		SELECT
			i.id AS inst_id,
			i.name AS inst_name,
			i.institution_code,
			i.state AS inst_state,
			i.is_active AS inst_active,

			d.id AS dept_id,
			d.department_name AS department_name,
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

		FROM institutions i

		LEFT JOIN departments d
			ON d.institution_id = i.id
			AND d.deleted_at IS NULL

		LEFT JOIN faculties f
			ON f.department_id = d.id
			AND f.deleted_at IS NULL

		LEFT JOIN students s
			ON s.faculty_id = f.id
			AND s.deleted_at IS NULL

		LEFT JOIN fees fe
			ON fe.student_id = s.id
			AND fe.deleted_at IS NULL

		WHERE i.id IN ?

	`, instIDs).Scan(&rows).Error

	if err != nil {
		return err
	}


	instIndex := make(map[uint]int)

	for i := range insts {
		instIndex[insts[i].ID] = i
	}


	deptIndex := make(map[uint]map[uint]int)
	facIndex := make(map[uint]map[uint]int)
	studIndex := make(map[uint]map[uint]int)


	for _, row := range rows {

		instPos, ok := instIndex[row.InstID]
		if !ok {
			continue
		}


		// Department
		if row.DeptID == nil {
			continue
		}


		if deptIndex[row.InstID] == nil {
			deptIndex[row.InstID] = make(map[uint]int)
		}


		deptPos, exists := deptIndex[row.InstID][*row.DeptID]

		if !exists {

			dept := model.Department{
				ID:            *row.DeptID,
				InstitutionID: row.InstID,
			}


			// Change DepartmentName to your actual model field
			if row.DepartmentName != nil {
				dept.DepartmentName = *row.DepartmentName
			}

			if row.DeptActive != nil {
				dept.IsActive = *row.DeptActive
			}


			insts[instPos].Departments =
				append(insts[instPos].Departments, dept)


			deptPos = len(insts[instPos].Departments) - 1

			deptIndex[row.InstID][*row.DeptID] = deptPos
		}



		// Faculty
		if row.FacID == nil {
			continue
		}


		dept := &insts[instPos].Departments[deptPos]


		if facIndex[*row.DeptID] == nil {
			facIndex[*row.DeptID] = make(map[uint]int)
		}


		facPos, exists := facIndex[*row.DeptID][*row.FacID]


		if !exists {

			fac := model.Faculty{
				ID:           *row.FacID,
				DepartmentID: *row.DeptID,
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


			dept.Faculties =
				append(dept.Faculties, fac)


			facPos = len(dept.Faculties) - 1

			facIndex[*row.DeptID][*row.FacID] = facPos
		}



		// Student
		if row.StudID == nil {
			continue
		}


		faculty := &dept.Faculties[facPos]


		if studIndex[*row.FacID] == nil {
			studIndex[*row.FacID] = make(map[uint]int)
		}


		studPos, exists := studIndex[*row.FacID][*row.StudID]


		if !exists {

			stud := model.Student{
				ID:        *row.StudID,
				FacultyID: *row.FacID,
			}


			if row.StudName != nil {
				stud.Name = *row.StudName
			}

			if row.StudEmail != nil {
				stud.Email = *row.StudEmail
			}

			if row.StudGender != nil {
				stud.Gender = *row.StudGender
			}

			if row.StudActive != nil {
				stud.IsActive = *row.StudActive
			}


			faculty.Students =
				append(faculty.Students, stud)


			studPos = len(faculty.Students) - 1

			studIndex[*row.FacID][*row.StudID] = studPos
		}



		// Fees
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
	if err != nil {
		return err
	}

	institute.ID = uint(id)
	institute.CreatedAt = now
	institute.UpdatedAt = now
	institute.IsActive = true

	if len(institute.Departments) > 0 {
		for i := range institute.Departments {
			institute.Departments[i].InstitutionID = institute.ID
			institute.Departments[i].CreatedAt = now
			institute.Departments[i].UpdatedAt = now
			institute.Departments[i].IsActive = true

			deptRes, deptErr := db.Exec(
				`INSERT INTO departments (department_name, institution_id, created_at, updated_at, is_active)
				VALUES (?, ?, ?, ?, ?)`,
				institute.Departments[i].DepartmentName,
				institute.Departments[i].InstitutionID,
				institute.Departments[i].CreatedAt,
				institute.Departments[i].UpdatedAt,
				institute.Departments[i].IsActive,
			)
			if deptErr != nil {
				return deptErr
			}
			deptID, deptIDErr := deptRes.LastInsertId()
			if deptIDErr != nil {
				return deptIDErr
			}
			institute.Departments[i].ID = uint(deptID)
		}
	}

	return nil
}

func (r *InstitutionRepository) FetchInstitution() ([]model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE deleted_at IS NULL").Scan(&insts).Error
	if err != nil {
		return nil, err
	}
	err = r.loadAssociations(insts)
	return insts, err
}

func (r *InstitutionRepository) FetchInstitutionPaginated(page, limit int) ([]model.Institutions, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM institutions WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var insts []model.Institutions
	err = r.db.Raw("SELECT * FROM institutions WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset).Scan(&insts).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.loadAssociations(insts)
	return insts, total, err
}

func (r *InstitutionRepository) FetchInstitutionById(id uint) (model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&insts).Error
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(insts)
	if err != nil {
		return model.Institutions{}, err
	}
	return insts[0], nil
}

func (r *InstitutionRepository) GetActiveInstitute() (model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true).Scan(&insts).Error
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(insts)
	if err != nil {
		return model.Institutions{}, err
	}
	return insts[0], nil
}

func (r *InstitutionRepository) GetInactiveInstitute() (model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false).Scan(&insts).Error
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	err = r.loadAssociations(insts)
	if err != nil {
		return model.Institutions{}, err
	}
	return insts[0], nil
}

func (r *InstitutionRepository) FetchInstitutionDeleted() ([]model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE deleted_at IS NOT NULL").Scan(&insts).Error
	if err != nil {
		return nil, err
	}
	err = r.loadAssociations(insts)
	return insts, err
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
