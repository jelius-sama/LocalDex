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
	"gopkg.in/gomail.v2"
	"net/http"
	"os"
	"strconv"
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

type SMTP struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Mail struct {
	From    string
	To      string
	Subject string
	Body    string
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

func sendEmail(smtp SMTP, mail Mail) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)

	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User, smtp.Password)
	d.SSL = false // iCloud uses STARTTLS, not implicit SSL
	return d.DialAndSend(m)
}

func buildEmailBody(otp string, expMinutes int) string {
	return fmt.Sprintf(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <style>
    @media (prefers-color-scheme: light) {
      .bg { background-color:#0b0d10 !important; }
    }
    a { text-decoration: none; }
  </style>
</head>
<body style="margin:0;padding:0;background-color:#050607;font-family:Inter,-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;">
  <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="min-width:320px;">
    <tr>
      <td align="center" style="padding:28px 16px;">
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="680" style="max-width:680px;width:100%%;border-radius:12px;background:linear-gradient(180deg,rgba(255,255,255,0.02),rgba(255,255,255,0.01));box-shadow:0 6px 28px rgba(2,6,23,0.6);overflow:hidden;">
          <tr>
            <td style="padding:28px 32px 18px 32px;color:#e6eef8;">
              <table role="presentation" width="100%%">
                <tr>
                  <td style="vertical-align:middle;">
                    <div style="display:flex;align-items:center;gap:12px;">
                      <div style="width:48px;height:48px;aspect-ratio:1/1;border-radius:10px;display:flex;align-items:center;justify-content:center;">
                        <img src="https://nas.jelius.dev/assets/favicon.png" style="width:100%%;height:100%%;object-fit:contain;aspect-ratio:1/1;" />
                      </div>
                      <div>
                        <div style="font-size:16px;font-weight:600;color:#ffffff;">LocalDex</div>
                        <div style="font-size:13px;color:#9fb0c8;margin-top:2px;">One-time passcode (secure)</div>
                      </div>
                    </div>
                  </td>
                  <td align="right" style="vertical-align:middle;color:#9fb0c8;font-size:13px;">
                    Expires in %d minutes
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:18px 32px 28px 32px;">
              <div style="font-size:15px;line-height:1.45;color:#cfe7ff;margin-bottom:18px;">
                Use the code below to sign in to your LocalDex account. This code can only be used once.
              </div>
              <div style="text-align:center;margin-bottom:22px;">
                <div style="display:inline-block;padding:18px 26px;border-radius:12px;background:linear-gradient(180deg,rgba(255,255,255,0.03),rgba(255,255,255,0.01));box-shadow:inset 0 -2px 8px rgba(0,0,0,0.45);">
                  <div style="font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,'Roboto Mono',monospace;letter-spacing:3px;font-size:28px;font-weight:700;color:#f8fafc;">
                    %s
                  </div>
                </div>
              </div>
              <div style="text-align:center;color:#9fb0c8;font-size:13px;margin-bottom:16px;">
                Expires in <strong style="color:#dfeefc;">%d minutes</strong> • Don’t share this with anyone
              </div>
              <div style="text-align:center;margin-bottom:6px;">
		<a href="https://nas.jelius.dev" style="display:inline-block;padding:10px 18px;border-radius:10px;border:1px solid rgba(255,255,255,0.06);background:linear-gradient(180deg,#0f1724,#071223);font-weight:600;color:#9ae6b4;font-size:14px;">
                  Open LocalDex
                </a>
              </div>
              <div style="border-top:1px solid rgba(255,255,255,0.02);color:#92b1cc;font-size:13px;padding-top:18px;line-height:1.45;">
                If you didn't request this, don't ignore this email becase someone might have tried to access my personal NAS.
              </div>
            </td>
          </tr>
          <tr>
            <td style="padding:18px 32px;background:#030405;color:#7f9db3;font-size:12px;text-align:center;">
              LocalDex · College Road · Bijni · Assam<br>
              <span style="color:#5f7d93;">For help, reply to this email.</span>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, expMinutes, otp, expMinutes)
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

	subject := "Your LocalDex OTP"
	body := buildEmailBody(otp, int(otpValidity/time.Minute))
	port, err := strconv.Atoi(smtpPort)

	if err != nil {
		logger.TimedError("Failed to parse port integer", "\n    ", err.Error())
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	smptProvider := SMTP{
		Host:     smtpHost,
		Port:     port,
		User:     smtpUser,
		Password: smtpPass,
	}

	if err := sendEmail(smptProvider, Mail{
		From:    "work@jelius.dev",
		To:      "personal@jelius.dev",
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
