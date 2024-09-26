package models

type PasswordUpdate struct {
	OldPassword       string `json:"old_password" example:"password"`
	NewPassword       string `json:"new_password" example:"new_password"`
	NewPasswordRepeat string `json:"new_password_repeat" example:"new_password"`
}
