package auth

import (
	"LocalDex/logger"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// OTP / session lifetimes
const (
	otpValidity     = 5 * time.Minute
	sessionValidity = 3 * time.Hour
)

// OTP store
var otpStore = struct {
	sync.RWMutex
	m map[string]otpEntry
}{m: make(map[string]otpEntry)}

// Session token store
var authTokens = struct {
	sync.RWMutex
	m map[string]time.Time
}{m: make(map[string]time.Time)}

type otpEntry struct {
	Code      string
	ExpiresAt time.Time
}

func init() {
	// GC for OTPs
	go func() {
		ticker := time.NewTicker(otpValidity)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			otpStore.Lock()
			for email, e := range otpStore.m {
				if now.After(e.ExpiresAt) {
					delete(otpStore.m, email)
				}
			}
			otpStore.Unlock()
			logger.TimedInfo("Expired OTPs cleaned by OTP GC.")
		}
	}()

	// GC for sessions
	go func() {
		ticker := time.NewTicker(sessionValidity)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			authTokens.Lock()
			for tok, exp := range authTokens.m {
				if now.After(exp) {
					delete(authTokens.m, tok)
				}
			}
			authTokens.Unlock()
			logger.TimedInfo("Expired sessions cleaned by AuthToken GC.")
		}
	}()
}

// generateSecureToken returns a securely generated random hex string of length 2*nBytes.
func generateSecureToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// verifyPassword checks the candidate against adminPasswordHash using salt+pepper.
func verifyPassword(candidate string) bool {
	if len(os.Getenv("PASSWORD_SALT")) == 0 {
		logger.TimedError("password salt isn't set properly")
		return false
	}
	if len(os.Getenv("PASSWORD_PEPPER")) == 0 {
		logger.TimedError("password pepper isn't set properly")
		return false
	}
	if len(os.Getenv("ADMIN_PASSWORD_HASH")) == 0 {
		logger.TimedError("admin password hash isn't set properly")
		return false
	}

	mac := hmac.New(sha256.New, []byte(os.Getenv("PASSWORD_PEPPER")))
	mac.Write([]byte(candidate + os.Getenv("PASSWORD_SALT")))
	expectedMAC := mac.Sum(nil)

	storedMAC, err := hex.DecodeString(os.Getenv("ADMIN_PASSWORD_HASH"))
	if err != nil {
		// misconfigured hash
		return false
	}

	// constant-time compare
	return subtle.ConstantTimeCompare(expectedMAC, storedMAC) == 1
}

// SendOTPHandler validates email+password, generates an OTP, emails it, and stores it.
func SendOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if len(os.Getenv("ADMIN_EMAIL")) == 0 {
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}

	// verify admin credentials
	if req.Email != os.Getenv("ADMIN_EMAIL") || !verifyPassword(req.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	otp, err := generateSecureToken(6) // 12-char hex
	if err != nil {
		http.Error(w, "failed to generate OTP", http.StatusInternalServerError)
		return
	}
	var (
		smtpHost = os.Getenv("SMTP_HOST")
		smtpPort = os.Getenv("SMTP_PORT")
		smtpUser = os.Getenv("SMTP_USER")
		smtpPass = os.Getenv("SMTP_PASS")
	)
	if len(smtpHost) == 0 || len(smtpPort) == 0 || len(smtpUser) == 0 || len(smtpPass) == 0 {
		logger.TimedError("smpt not set up correctly")
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}

	// store OTP
	otpStore.Lock()
	otpStore.m[req.Email] = otpEntry{
		Code:      otp,
		ExpiresAt: time.Now().Add(otpValidity),
	}
	otpStore.Unlock()

	// send email using STARTTLS
	subject := "Your LocalDex OTP"
	body := fmt.Sprintf("Your LocalDex login OTP is: %s\nIt expires in %d minutes.", otp, int(otpValidity/time.Minute))

	if err := sendMailSMTP587(SendMail{
		Host:    smtpHost,
		Port:    smtpPort,
		User:    smtpUser,
		Pass:    smtpPass,
		From:    "work@jelius.dev",
		To:      req.Email,
		Subject: subject,
		Body:    body,
	}); err != nil {
		logger.TimedInfo("failed to send OTP email: " + err.Error())
		http.Error(w, "failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"OTP sent"}`))
}

// VerifyOTPHandler checks the OTP and, if valid, issues a session cookie.
func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// verify OTP
	otpStore.RLock()
	entry, ok := otpStore.m[req.Email]
	otpStore.RUnlock()
	if !ok || time.Now().After(entry.ExpiresAt) || entry.Code != req.OTP {
		http.Error(w, "invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// consume OTP
	otpStore.Lock()
	delete(otpStore.m, req.Email)
	otpStore.Unlock()

	// issue session token
	token, err := generateSecureToken(32) // 64-char hex
	if err != nil {
		http.Error(w, "failed to generate session token", http.StatusInternalServerError)
		return
	}
	expiry := time.Now().Add(sessionValidity)

	authTokens.Lock()
	authTokens.m[token] = expiry
	authTokens.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"authenticated"}`))
}

// VerifyAuthToken checks and renews a session token.
func VerifyAuthToken(token string) error {
	authTokens.RLock()
	exp, ok := authTokens.m[token]
	authTokens.RUnlock()
	if !ok {
		return fmt.Errorf("token not found")
	}
	if time.Now().After(exp) {
		authTokens.Lock()
		delete(authTokens.m, token)
		authTokens.Unlock()
		return fmt.Errorf("token expired")
	}

	// renew
	newExp := time.Now().Add(sessionValidity)
	authTokens.Lock()
	authTokens.m[token] = newExp
	authTokens.Unlock()
	return nil
}

// VerifyAuthStatus is a protected handler to check/renew the session cookie.
func VerifyAuthStatus(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "missing auth token", http.StatusForbidden)
		return
	}

	if err := VerifyAuthToken(c.Value); err != nil {
		if err.Error() == "token expired" {
			http.Error(w, err.Error(), 498)
		} else {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
		return
	}

	// renew cookie
	newExpiry := time.Now().Add(sessionValidity)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    c.Value,
		Expires:  newExpiry,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("authorized"))
}
