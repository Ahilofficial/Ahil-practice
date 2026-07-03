package dto

import "time"

type CreateFacultyDTO struct {
	Name         string    `json:"name"`
	Gender       string    `json:"gender"`
	JoiningDate  time.Time `json:"joining_date"`
	DepartmentID uint      `json:"department_id"`
}
