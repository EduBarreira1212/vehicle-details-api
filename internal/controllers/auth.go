package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/EduBarreira1212/vehicle-details-api/internal/auth"
	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"github.com/EduBarreira1212/vehicle-details-api/internal/repositories"
	"github.com/EduBarreira1212/vehicle-details-api/internal/responses"
	"github.com/EduBarreira1212/vehicle-details-api/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordInput struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func ForgotPassword(c *gin.Context) {
	var in ForgotPasswordInput
	if err := c.ShouldBindJSON(&in); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	successMsg := gin.H{"message": "Se o e-mail existir, enviaremos instruções para redefinir a senha."}

	repository := repositories.NewUserRepository(config.DB)
	user, err := repository.GetByEmail(c.Request.Context(), in.Email)
	if err != nil {
		responses.JSON(c.Writer, http.StatusOK, successMsg)
		return
	}

	token, raw, err := auth.GenerateSecureToken(32)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	tokenHash := auth.HashToken(raw)

	tokenRepository := repositories.NewTokenRepository(config.DB)
	if err := tokenRepository.ExpireActiveResetTokens(c.Request.Context(), user.ID); err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	reset := models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := tokenRepository.CreatePasswordResetToken(c.Request.Context(), &reset); err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	//ENVIAR E-MAIL
	//resetURL := appBaseURL + "/reset-password?token=" + token

	//_ = mailer.Send(user.Email, "Redefinição de senha",
	//	"<p>Para redefinir sua senha, clique no link abaixo:</p>"+
	//		`<p><a href="`+resetURL+`">Redefinir senha</a></p>`+
	//		"<p>Se você não solicitou, ignore este e-mail.</p>")

	responses.JSON(c.Writer, http.StatusOK, successMsg)
}

func ResetPassword(c *gin.Context) {
	var in ResetPasswordInput
	if err := c.ShouldBindJSON(&in); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}
	if len(in.Password) < 8 {
		responses.Error(c.Writer, http.StatusBadRequest, errors.New("password needs to be above 8 characters"))
		return
	}

	raw, err := utils.DecodeBase64URL(in.Token)
	if err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}
	tokenHash := auth.HashToken(raw)

	tokenRepository := repositories.NewTokenRepository(config.DB)
	prt, err := tokenRepository.FindPasswordResetTokenByHash(c.Request.Context(), string(tokenHash))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.Error(c.Writer, http.StatusBadRequest, err)
			return
		}

		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	if prt.UsedAt != nil || time.Now().After(prt.ExpiresAt) {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	userRepository := repositories.NewUserRepository(config.DB)
	user, err := userRepository.GetById(c.Request.Context(), prt.ID)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	hashed, err := auth.Hash(in.Password)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	if err := tokenRepository.ResetPasswordWithToken(c.Request.Context(), user.ID, prt.ID, string(hashed)); err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusOK, gin.H{"message": "Senha redefinida com sucesso."})
}
