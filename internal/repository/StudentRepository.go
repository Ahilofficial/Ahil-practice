package repository

import (
	"backend_institutions/internal/model"
	"errors"
	"time"
	"backend_institutions/internal/dto"
	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}


func (r *StudentRepository) fetchWithRelations(baseQuery string, args ...interface{}) ([]model.Student, error) {
	var rows []dto.StudentFlatRow
	query := `
	SELECT 
		s.id AS stud_id, s.name AS stud_name, s.email AS stud_email, s.gender AS stud_gender, s.faculty_id, s.is_active AS stud_active,
		fe.id AS fee_id, fe.payment_mode AS fee_payment_mode, fe.amount AS fee_amount, fe.is_active AS fee_active
	FROM (` + baseQuery + `) s
	LEFT JOIN fees fe ON fe.student_id = s.id AND fe.deleted_at IS NULL
	ORDER BY s.id, fe.id`

	err := r.db.Raw(query, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	studMap := make(map[uint]*model.Student)
	var orderedIDs []uint
	for _, row := range rows {
		stud, exists := studMap[row.StudID]
		if !exists {
			stud = &model.Student{
				ID:        row.StudID,
				Name:      row.StudName,
				Email:     row.StudEmail,
				Gender:    row.StudGender,
				FacultyID: row.FacultyID,
				IsActive:  row.StudActive,
				Fees:      []model.Fees{},
			}
			studMap[row.StudID] = stud
			orderedIDs = append(orderedIDs, row.StudID)
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
					StudentID:   row.StudID,
					IsActive:    *row.FeeActive,
				}
				stud.Fees = append(stud.Fees, fee)
			}
		}
	}

	result := make([]model.Student, len(orderedIDs))
	for i, id := range orderedIDs {
		result[i] = *studMap[id]
	}
	return result, nil
}

func (r *StudentRepository) CreateStudent(student *model.Student) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	res, err := db.Exec(
		`INSERT INTO students (name, email, gender, faculty_id, created_at, updated_at, is_active)
		SELECT ?, ?, ?, id, ?, ?, ? FROM faculties
		WHERE id = ? AND deleted_at IS NULL AND is_active = true
		  AND NOT EXISTS (
			  SELECT 1 FROM students 
			  WHERE email = ? AND deleted_at IS NULL
		  )`,
		student.Name, student.Email, student.Gender, now, now, true,
		student.FacultyID,
		student.Email,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("student email already registered, or assigned faculty is inactive/invalid")
	}
	id, err := res.LastInsertId()
	if err == nil {
		student.ID = uint(id)
	}
	return nil
}

func (r *StudentRepository) FetchStudent() ([]model.Student, error) {
	return r.fetchWithRelations("SELECT * FROM students WHERE deleted_at IS NULL")
}

func (r *StudentRepository) GetActiveStudent() (model.Student, error) {
	studs, err := r.fetchWithRelations("SELECT * FROM students WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true)
	if err != nil {
		return model.Student{}, err
	}
	if len(studs) == 0 {
		return model.Student{}, gorm.ErrRecordNotFound
	}
	return studs[0], nil
}

func (r *StudentRepository) GetInactiveStudent() (model.Student, error) {
	studs, err := r.fetchWithRelations("SELECT * FROM students WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false)
	if err != nil {
		return model.Student{}, err
	}
	if len(studs) == 0 {
		return model.Student{}, gorm.ErrRecordNotFound
	}
	return studs[0], nil
}

func (r *StudentRepository) FetchStudentPaginated(page, limit int) ([]model.Student, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM students WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	studs, err := r.fetchWithRelations("SELECT * FROM students WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset)
	return studs, total, err
}

func (r *StudentRepository) FetchStudentById(id uint) (model.Student, error) {
	studs, err := r.fetchWithRelations("SELECT * FROM students WHERE id = ? AND deleted_at IS NULL LIMIT 1", id)
	if err != nil {
		return model.Student{}, err
	}
	if len(studs) == 0 {
		return model.Student{}, gorm.ErrRecordNotFound
	}
	return studs[0], nil
}

func (r *StudentRepository) FetchStudentDeleted() ([]model.Student, error) {
	return r.fetchWithRelations("SELECT * FROM students WHERE deleted_at IS NOT NULL")
}

func (r *StudentRepository) DeleteStudent(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE students SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
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

func (r *StudentRepository) UpdateStudentById(student *model.Student) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE students SET name = ?, email = ?, gender = ?, updated_at = ? WHERE id = ?",
		student.Name, student.Email, student.Gender, time.Now(), student.ID,
	)
	return err
}
