package repositories

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"warkop-api/dto"
	"warkop-api/helpers"
	"warkop-api/models"

	"gorm.io/gorm"
)

func (r *compRepository) RegisterUser(data dto.User) (*string, error) {
	user := models.User{
		Username:  data.Username,
		Email:     data.Email,
		Password:  data.Password,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Contact:   data.Contact,
		Address:   data.Address,
	}

	err := r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (r *compRepository) RegisterToken(data dto.User) (*string, error) {
	token, err := helpers.GenerateToken(32)
	if err != nil {
		return nil, err
	}

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", data.ID).Delete(&models.VerificationToken{}).Error
		if err != nil {
			return err
		}

		verificationToken := models.VerificationToken{
			UserID:    data.ID,
			Token:     token,
			ExpiredAt: time.Now().Add(9 * time.Hour),
		}

		return tx.Create(&verificationToken).Error
	})
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *compRepository) VerifyAccount(token string) error {
	var verificationToken models.VerificationToken

	err := r.DB.Where("token = ?", token).First(&verificationToken).Error
	if err != nil {
		return err
	}

	if time.Now().After(verificationToken.ExpiredAt) {
		return errors.New(strconv.Itoa(http.StatusGone))
	}

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		err := tx.Model(&models.User{}).Where("id = ?", verificationToken.UserID).Updates(map[string]interface{}{
			"is_verified": true,
			"verified_at": &now,
		}).Error
		if err != nil {
			return err
		}

		return tx.Where("token = ?", token).Delete(&models.VerificationToken{}).Error
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *compRepository) GetUserData(username string) (*dto.User, error) {
	var user models.User

	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	var verifiedAt *string
	if user.VerifiedAt != nil {
		vStr := user.VerifiedAt.Format("2006-01-02 15:04:05")
		verifiedAt = &vStr
	}

	data := dto.User{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Password:   user.Password,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Contact:    user.Contact,
		Address:    user.Address,
		IsVerified: user.IsVerified,
		VerifiedAt: verifiedAt,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &data, nil
}

func (r *compRepository) RequestResetPassword(data dto.User, otp string) error {
	otpInt, err := strconv.Atoi(otp)
	if err != nil {
		return err
	}

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", data.ID).Delete(&models.ResetOTP{}).Error
		if err != nil {
			return err
		}

		resetOTP := models.ResetOTP{
			UserID:    data.ID,
			OTP:       otpInt,
			ExpiredAt: time.Now().Add(2 * time.Hour),
		}

		return tx.Create(&resetOTP).Error
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *compRepository) VerifyResetPassword(data dto.OTPVerifyToken) (*dto.OTPVerifyToken, error) {
	var resetOTP models.ResetOTP

	err := r.DB.Where("user_id = ?", data.UserID).First(&resetOTP).Error
	if err != nil {
		return nil, err
	}

	d := dto.OTPVerifyToken{
		UserID: resetOTP.UserID,
		OTP:    strconv.Itoa(resetOTP.OTP),
	}

	return &d, nil
}

func (r *compRepository) ResetPassword(user_data dto.User) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", user_data.ID).Delete(&models.ResetOTP{}).Error
		if err != nil {
			return err
		}

		return tx.Model(&models.User{}).Where("id = ?", user_data.ID).Update("password", user_data.Password).Error
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *compRepository) UploadUserProfile(data dto.User, image_url string) error {
	var userImage models.UsersImage
	err := r.DB.Where("user_id = ?", data.ID).First(&userImage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userImage = models.UsersImage{
				UserID:   data.ID,
				ImageURL: image_url,
			}
			return r.DB.Create(&userImage).Error
		}
		return err
	}

	return r.DB.Model(&userImage).Update("image_url", image_url).Error
}

func (r *compRepository) GetUserProfile(id string) (*string, error) {
	var userImage models.UsersImage

	err := r.DB.Where("user_id = ?", id).First(&userImage).Error
	if err != nil {
		return nil, err
	}

	return &userImage.ImageURL, nil
}
