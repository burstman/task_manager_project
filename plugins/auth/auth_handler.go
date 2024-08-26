package auth

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"taskManager/app/db"
	"time"

	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	UserSessionName = "user-session"
)

var authSchema = v.Schema{
	"email":    v.Rules(v.Email),
	"password": v.Rules(v.Required),
}

func HandleLoginIndex(kit *kit.Kit) error {
	if kit.Auth().Check() {
		redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")
		return kit.Redirect(http.StatusSeeOther, redirectURL)
	}
	return kit.Render(LoginIndex(LoginIndexPageData{}))
}

// HandleLoginCreate handles the login form submission. It validates the form data,
// checks the user's credentials, and creates a new session if the login is successful.
// If the user's email is not verified, it adds an error to the form and re-renders the login form.
// If the login is successful, it redirects the user to the configured redirect URL after login.
func HandleLoginCreate(kit *kit.Kit) error {
	var values LoginFormValues
	errors, ok := v.Request(kit.Request, &values, authSchema)
	if !ok {
		return kit.Render(LoginForm(values, errors))
	}

	var user User
	err := db.Get().Find(&user, "email = ?", values.Email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.Add("credentials", "invalid credentials")
			return kit.Render(LoginForm(values, errors))
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(values.Password))
	if err != nil {
		errors.Add("credentials", "invalid credentials")
		return kit.Render(LoginForm(values, errors))
	}

	skipVerify := kit.Getenv("SUPERKIT_AUTH_SKIP_VERIFY", "false")
	if skipVerify != "true" {
		if !user.EmailVerifiedAt.Valid {
			errors.Add("verified", "please verify your email")
			return kit.Render(LoginForm(values, errors))
		}
	}

	sessionExpiryStr := kit.Getenv("SUPERKIT_AUTH_SESSION_EXPIRY_IN_HOURS", "48")
	sessionExpiry, err := strconv.Atoi(sessionExpiryStr)
	if err != nil {
		sessionExpiry = 48
	}
	session := Session{
		User:      user,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(sessionExpiry)),
	}
	if err = db.Get().Create(&session).Error; err != nil {
		return err
	}

	sess := kit.GetSession(UserSessionName)
	sess.Values["sessionToken"] = session.Token
	sess.Values["firstName"] = session.User.FirstName
	sess.Values["userid"] = session.User.ID
	sess.Save(kit.Request, kit.Response)
	redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")

	return kit.Redirect(http.StatusSeeOther, redirectURL)
}

func HandleLoginDelete(kit *kit.Kit) error {
	sess := kit.GetSession(UserSessionName)
	defer func() {
		sess.Values = map[any]any{}
		sess.Save(kit.Request, kit.Response)
	}()
	err := db.Get().Delete(&Session{}, "token = ?", sess.Values["sessionToken"]).Error
	if err != nil {
		return err
	}
	return kit.Redirect(http.StatusSeeOther, "/")
}

func HandleEmailVerify(kit *kit.Kit) error {
	tokenStr := kit.Request.URL.Query().Get("token")
	if len(tokenStr) == 0 {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	token, err := jwt.ParseWithClaims(
		tokenStr, &jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SUPERKIT_SECRET")), nil
		}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if !token.Valid {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	var user User
	err = db.Get().First(&user, userID).Error
	if err != nil {
		return err
	}

	if user.EmailVerifiedAt.Time.After(time.Time{}) {
		return kit.Render(EmailVerificationError("Email already verified"))
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	user.EmailVerifiedAt = now
	err = db.Get().Save(&user).Error
	if err != nil {
		return err
	}

	return kit.Redirect(http.StatusSeeOther, "/login")
}

// AuthenticateUser retrieves the user's session token from the session and checks if the token is valid and not expired.
// If the token is valid, it returns an Auth struct with the user's ID, email, and a LoggedIn flag set to true.
// If the token is invalid or expired, it returns an empty Auth struct.
func AuthenticateUser(kit *kit.Kit) (kit.Auth, error) {
	auth := Auth{}
	sess := kit.GetSession(UserSessionName)
	token, ok := sess.Values["sessionToken"]
	//fmt.Println("ok:", ok)
	if !ok {
		return auth, nil
	}

	var session Session
	err := db.Get().
		Preload("User").
		Find(&session, "token = ? AND expires_at > ?", token, time.Now()).Error
	if err != nil || session.ID == 0 {
		return auth, nil
	}

	return Auth{
		LoggedIn: true,
		UserID:   session.User.ID,
		Email:    session.User.Email,
	}, nil
}

func UserHasRole(userID int, roleName string) bool {
	// Query the database to check if the user has the specified role
	return true
}

func UserHasPermission(userID uuid.UUID, permissionName string) bool {
	// Query the database to check if the user has the specified permission
	return true
}
