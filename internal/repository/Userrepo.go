package repository

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"errors"
	"fmt"

	// "fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByVerificationToken(token string) (model.User, error) {
	var user model.User

	query := `
		SELECT *
		FROM users
		WHERE verification_token = ?
		LIMIT 1
	`

	err := r.db.Raw(query, token).Scan(&user).Error
	if err != nil {
		return model.User{}, err
	}
	if user.ID == 0 {
		return model.User{}, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	var expiresAt any = user.TokenExpiresAt
	if user.TokenExpiresAt.IsZero() {
		expiresAt = nil
	}

	query := `
		UPDATE users
		SET
			is_active = ?,
			is_verified = ?,
			verification_token = ?,
			token_expires_at = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	return r.db.Exec(
		query,
		user.IsActive,
		user.IsVerified,
		user.VerificationToken,
		expiresAt,
		user.ID,
	).Error
}

func (r *UserRepository) CreateUser(user *model.User) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	var expiresAt any = user.TokenExpiresAt
	if user.TokenExpiresAt.IsZero() {
		expiresAt = nil
	}
	query := `
		INSERT INTO users (name, email, phone, password, is_active, is_verified, verification_token, token_expires_at, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := db.Exec(
		query,
		user.Name, user.Email, user.Phone, user.Password, user.IsActive, user.IsVerified, user.VerificationToken, expiresAt, now, now,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		user.ID = uint(id)
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Raw("SELECT * FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1", email).Scan(&user).Error
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (r *UserRepository) FindByPhone(phone string) (model.User, error) {
	var user model.User
	err := r.db.Raw("SELECT * FROM users WHERE phone = ? AND deleted_at IS NULL LIMIT 1", phone).Scan(&user).Error
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (r *UserRepository) AssignRoleToUser(userID uint, roleID uint) error {
	if err := r.db.Exec("DELETE FROM user_roles WHERE user_id = ?", userID).Error; err != nil {
		return err
	}

	result := r.db.Exec(
		"INSERT INTO user_roles (user_id, role_id) SELECT ?, id FROM roles WHERE id = ?",
		userID,
		roleID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("role does not exist")
	}

	return nil
}

func (r *UserRepository) FindRoleByName(name string) (model.Role, error) {
	var role model.Role
	err := r.db.Raw("SELECT id, name FROM roles WHERE name = ? LIMIT 1", name).Scan(&role).Error
	if err != nil {
		return role, err
	}
	if role.ID == 0 {
		return role, gorm.ErrRecordNotFound
	}
	return role, nil
}

func (r *UserRepository) DeleteUser(id uint) error {

	if err := r.db.Exec("DELETE FROM user_roles WHERE user_id = ?", id).Error; err != nil {
		return err
	}

	res := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}


func(s *UserRepository)ForgotPasswordRepo(dto dto.ForgotPasswordDTO)(model.User,error){
	var user model.User
	query:=`select * from users where email=? limit 1`
	result:= s.db.Raw(query,dto.Email).Scan(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}


	if result.RowsAffected == 0 {
		return model.User{}, errors.New("email not found")
	}

	return user,nil
}

func (r *UserRepository) UpdateResetToken(user model.User) error {
	query := `
		UPDATE users
		SET
			reset_password_token = ?,
			reset_token_expires_at = ?
		WHERE id = ?
	`

	return r.db.Exec(
		query,
		user.ResetPasswordToken,
		user.ResetTokenExpiresAt,
		user.ID,
	).Error
}

func (r *UserRepository) FetchUsertoken(token string) (model.User, error) {
    var user model.User

    fmt.Println("Searching Token:", token)

    query := `
        SELECT *
        FROM users
        WHERE reset_password_token = ?
        LIMIT 1
    `

    result := r.db.Raw(query, token).Scan(&user)

    fmt.Println("Rows Affected:", result.RowsAffected)
    fmt.Println("DB Error:", result.Error)
    fmt.Printf("User: %+v\n", user)

    if result.Error != nil {
        return model.User{}, result.Error
    }

    if result.RowsAffected == 0 {
        return model.User{}, errors.New("invalid reset token")
    }

    return user, nil
}

func (r *UserRepository) UpdatePassword(id uint, password string) error {
	
	query := `
		UPDATE users
		SET
			password = ?,
			reset_password_token = NULL,
			reset_token_expires_at = NULL
		WHERE id = ?
	`

	return r.db.Exec(query, password, id).Error
}


func (r *UserRepository) Logout(dto *dto.LogoutDTO) error {
	var sessionID string
	query := `
		SELECT session_id
		FROM sessions
		WHERE user_id = ?
		  AND refresh_token = ?
		  AND is_active = TRUE
		LIMIT 1
	`

	err := r.db.Raw(query, dto.UserID, dto.Token).Scan(&sessionID).Error
	if err != nil {
		return err
	}

	if sessionID == "" {
		return errors.New("invalid session or already logged out")
	}

	update := `
		UPDATE sessions
		SET
			access_token = NULL,
			refresh_token = NULL,
			is_active = FALSE
		WHERE
			session_id = ?
	`

	return r.db.Exec(update, sessionID).Error
}

	func (r *UserRepository) FindByID(userID uint) (model.User, error) {
	var user model.User

	err := r.db.
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}