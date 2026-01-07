package utils

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net"
	"net/smtp"
	"os"
)

func GenerateVerificationCode() string {
	const chars = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[n.Int64()]
	}
	return string(code)
}

func SendVerificationEmail(to, code string) error {
	from := os.Getenv("SMTP_FROM")
	username := os.Getenv("SMTP_USERNAME")
	if username == "" {
		username = from
	}
	fromName := "Arizona Games Store"
	if envName := os.Getenv("SMTP_FROM_NAME"); envName != "" {
		fromName = envName
	}
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if smtpHost == "" || smtpPort == "" || from == "" || password == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	subject := "Arizona Games Store"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: Arial, sans-serif;
            background: linear-gradient(180deg, #000000 0%%, #1a0a0a 50%%, #2d0f1f 100%%);
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            padding: 20px;
        }
        .email-box {
            background: rgba(26, 10, 10, 0.9);
            border-radius: 15px;
            padding: 40px 30px;
            box-shadow: 0 10px 40px rgba(220, 20, 60, 0.3);
            border: 1px solid rgba(220, 20, 60, 0.2);
        }
        .logo {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo-text {
            font-size: 32px;
            font-weight: 700;
            color: #dc143c;
            text-shadow: 0 0 20px rgba(220, 20, 60, 0.5);
            margin: 0;
            letter-spacing: 2px;
        }
        .title {
            color: #ffffff;
            font-size: 24px;
            font-weight: 700;
            text-align: center;
            margin: 20px 0;
            text-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
        }
        .description {
            color: #888888;
            font-size: 16px;
            text-align: center;
            line-height: 1.6;
            margin: 20px 0 30px 0;
        }
        .code-container {
            background: linear-gradient(135deg, #dc143c 0%%, #b8102f 100%%);
            border-radius: 12px;
            padding: 30px;
            margin: 30px 0;
            box-shadow: 0 5px 20px rgba(220, 20, 60, 0.4);
            border: 2px solid rgba(255, 255, 255, 0.1);
        }
        .code-label {
            color: rgba(255, 255, 255, 0.8);
            font-size: 14px;
            text-align: center;
            margin-bottom: 15px;
            text-transform: uppercase;
            letter-spacing: 2px;
        }
        .code {
            background: rgba(0, 0, 0, 0.3);
            color: #ffffff;
            font-size: 42px;
            font-weight: 700;
            text-align: center;
            padding: 20px;
            border-radius: 8px;
            letter-spacing: 8px;
            font-family: 'Courier New', monospace;
            text-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
            border: 1px solid rgba(255, 255, 255, 0.2);
        }
        .timer {
            color: #ffcccc;
            font-size: 14px;
            text-align: center;
            margin-top: 15px;
        }
        .warning {
            background: rgba(220, 20, 60, 0.1);
            border-left: 4px solid #dc143c;
            padding: 15px 20px;
            margin: 25px 0;
            border-radius: 5px;
        }
        .warning-text {
            color: #ff6b6b;
            font-size: 14px;
            margin: 0;
            line-height: 1.5;
        }
        .footer {
            color: #555555;
            font-size: 13px;
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid rgba(220, 20, 60, 0.2);
            line-height: 1.6;
        }
        .footer-link {
            color: #dc143c;
            text-decoration: none;
        }
        @media only screen and (max-width: 600px) {
            .container {
                margin: 20px auto;
                padding: 10px;
            }
            .email-box {
                padding: 30px 20px;
            }
            .code {
                font-size: 32px;
                letter-spacing: 5px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="email-box">
            <div class="logo">
                <h1 class="logo-text">🎮 ARIZONA GAMES STORE</h1>
            </div>

            <h2 class="title">Добро пожаловать!</h2>

            <p class="description">
                Спасибо за регистрацию в игровом магазине Arizona Role Play!<br>
                Введите код ниже для подтверждения вашего аккаунта.
            </p>

            <div class="code-container">
                <div class="code-label">Ваш секретный код</div>
                <div class="code">%s</div>
                <div class="timer">⏱️ Код действителен 10 минут</div>
            </div>

            <div class="warning">
                <p class="warning-text">
                    ⚠️ Никому не сообщайте этот код! Администрация никогда не попросит вас отправить код подтверждения.
                </p>
            </div>

            <div class="footer">
                Если вы не регистрировались на нашем сайте, просто проигнорируйте это письмо.<br>
                <br>
                С уважением,<br>
                <strong style="color: #dc143c;">Команда Arizona Games Store</strong><br>
                <br>
                Arizona Role Play - Игровой магазин для проекта Arizona RP<br>
                Здесь вы можете покупать, продавать и арендовывать игровое имущество! 🎯
            </div>
        </div>
    </div>
</body>
</html>
`, code)

	message := []byte(
		"From: " + fromName + " <" + from + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body)

	testMode := os.Getenv("EMAIL_TEST_MODE")
	if testMode == "true" {
		fmt.Printf("📧 [TEST MODE] Email would be sent to: %s\n", to)
		fmt.Printf("📧 [TEST MODE] Code: %s\n", code)
		fmt.Printf("📧 [TEST MODE] Subject: %s\n", subject)
		return nil
	}

	serverAddr := net.JoinHostPort(smtpHost, smtpPort)

	conn, err := smtp.Dial(serverAddr)
	if err != nil {
		fmt.Printf("❌ Failed to connect: %v\n", err)
		return err
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = conn.StartTLS(tlsConfig); err != nil {
		fmt.Printf("❌ STARTTLS failed: %v\n", err)
		return err
	}
	fmt.Printf("✅ STARTTLS enabled\n")

	auth := smtp.PlainAuth("", username, password, smtpHost)
	if err = conn.Auth(auth); err != nil {
		fmt.Printf("❌ Auth failed: %v\n", err)
		return err
	}
	fmt.Printf("✅ Authentication successful\n")

	if err = conn.Mail(from); err != nil {
		fmt.Printf("❌ MAIL FROM failed: %v\n", err)
		return err
	}

	if err = conn.Rcpt(to); err != nil {
		fmt.Printf("❌ RCPT TO failed: %v\n", err)
		return err
	}

	w, err := conn.Data()
	if err != nil {
		fmt.Printf("❌ DATA command failed: %v\n", err)
		return err
	}

	_, err = w.Write(message)
	if err != nil {
		fmt.Printf("❌ Write failed: %v\n", err)
		return err
	}

	err = w.Close()
	if err != nil {
		fmt.Printf("❌ Close failed: %v\n", err)
		fmt.Printf("📧 [TEST MODE] Code for testing: %s\n", code)
		return err
	}

	conn.Quit()

	fmt.Printf("✅ Email sent successfully to %s\n", to)
	return nil
}
