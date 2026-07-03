package dto

type CreateStudentDTO struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	FacultyID uint   `json:"faculty_id"`
}
