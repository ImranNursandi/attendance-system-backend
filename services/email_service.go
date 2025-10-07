package services

import (
	"attendance-system/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ResendEmailService struct {
	apiKey    string
	fromEmail string
	config    *config.Config
}

type ResendSendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
	Text    string   `json:"text,omitempty"`
}

type ResendResponse struct {
	ID string `json:"id"`
}

type ResendError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewEmailService() *ResendEmailService {
	config := config.GetConfig()
	return &ResendEmailService{
		apiKey:    config.ResendAPIKey,
		fromEmail: config.FromEmail,
		config:    config,
	}
}

func (es *ResendEmailService) SendEmail(to, subject, body string) error {
	if es.apiKey == "" {
		// In development, log instead of failing
		fmt.Printf("📧 [DEV] Email would be sent to: %s\n", to)
		fmt.Printf("📧 [DEV] Subject: %s\n", subject)
		fmt.Printf("📧 [DEV] Body: %s\n", body)
		return nil
	}

	htmlBody := es.generateHTMLTemplate(subject, body)
	
	emailData := ResendSendRequest{
		From:    es.fromEmail,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
		Text:    body,
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		return fmt.Errorf("failed to marshal email data: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+es.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var resendErr ResendError
		if err := json.NewDecoder(resp.Body).Decode(&resendErr); err != nil {
			return fmt.Errorf("failed to send email: status %d", resp.StatusCode)
		}
		return fmt.Errorf("failed to send email: %s (status %d)", resendErr.Message, resendErr.Status)
	}

	var resendResp ResendResponse
	if err := json.NewDecoder(resp.Body).Decode(&resendResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf("✅ Email sent successfully to %s with ID: %s\n", to, resendResp.ID)
	return nil
}

func (es *ResendEmailService) generateHTMLTemplate(subject, body string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 30px 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
        }
        .content {
            padding: 30px;
        }
        .button {
            display: inline-block;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            text-decoration: none;
            padding: 12px 30px;
            border-radius: 5px;
            font-weight: 600;
            margin: 20px 0;
        }
        .footer {
            background-color: #f8f9fa;
            padding: 20px;
            text-align: center;
            color: #6c757d;
            font-size: 14px;
        }
        .code-block {
            background-color: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 4px;
            padding: 15px;
            margin: 15px 0;
            font-family: 'Courier New', monospace;
            word-break: break-all;
            font-size: 12px;
        }
        .info-box {
            background-color: #e7f3ff;
            border: 1px solid #b3d9ff;
            border-radius: 4px;
            padding: 15px;
            margin: 15px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Attendance System</h1>
        </div>
        <div class="content">
            %s
        </div>
        <div class="footer">
            <p>&copy; %d Company Name. All rights reserved.</p>
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>`, subject, body, time.Now().Year())
}

func (es *ResendEmailService) SendAccountSetupEmail(email, name, setupToken string) error {
	setupURL := fmt.Sprintf("%s/setup-account?token=%s", es.config.FrontendURL, setupToken)
	
	subject := "Set Up Your Attendance System Account"
	
	body := fmt.Sprintf(`
	<h2>Welcome to the Attendance System, %s!</h2>
	
	<div class="info-box">
		<strong>Your account has been created and is ready for setup.</strong>
	</div>
	
	<p>To activate your account and set your password, click the button below:</p>
	
	<div style="text-align: center;">
		<a href="%s" class="button">Set Up Your Account</a>
	</div>
	
	<p><strong>Can't click the button?</strong> Copy and paste this URL into your browser:</p>
	<div class="code-block">%s</div>
	
	<div class="info-box">
		<strong>Important:</strong> This link will expire in 7 days.
	</div>
	
	<p>If you didn't request this account, please contact your manager or ignore this email.</p>
	
	<p>Best regards,<br><strong>Attendance System Team</strong></p>
	`, name, setupURL, setupURL)

	return es.SendEmail(email, subject, body)
}

func (es *ResendEmailService) SendPasswordResetEmail(email, resetToken string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", es.config.FrontendURL, resetToken)
	
	subject := "Reset Your Attendance System Password"
	
	body := fmt.Sprintf(`
	<h2>Password Reset Request</h2>
	
	<div class="info-box">
		<strong>We received a request to reset your password.</strong>
	</div>
	
	<p>Click the button below to reset your password:</p>
	
	<div style="text-align: center;">
		<a href="%s" class="button">Reset Password</a>
	</div>
	
	<p><strong>Can't click the button?</strong> Copy and paste this URL into your browser:</p>
	<div class="code-block">%s</div>
	
	<div class="info-box">
		<strong>Important:</strong> This link will expire in 1 hour.
	</div>
	
	<p>If you didn't request a password reset, please ignore this email.</p>
	
	<p>Best regards,<br><strong>Attendance System Team</strong></p>
	`, resetURL, resetURL)

	return es.SendEmail(email, subject, body)
}

// SendWelcomeEmail sends a welcome email after account setup
func (es *ResendEmailService) SendWelcomeEmail(email, name string) error {
	subject := "Welcome to Attendance System"
	
	body := fmt.Sprintf(`
	<h2>Welcome aboard, %s! 🎉</h2>
	
	<div class="info-box">
		<strong>Your account has been successfully activated!</strong>
	</div>
	
	<p>You can now access the Attendance System using your credentials.</p>
	
	<p><strong>Quick Start Guide:</strong></p>
	<ul>
		<li>Clock in when you start work</li>
		<li>Clock out when you finish</li>
		<li>View your attendance history</li>
		<li>Update your profile information</li>
	</ul>
	
	<p>If you have any questions, please contact your manager.</p>
	
	<p>Best regards,<br><strong>Attendance System Team</strong></p>
	`, name)

	return es.SendEmail(email, subject, body)
}