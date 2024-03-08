package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/mail"
	"net/smtp"
	"sync"
	"time"

	"github.com/DmiTryAgain/secretSanta/pkg/calculator"

	"github.com/BurntSushi/toml"
)

const configFilename = "conf.toml"

type Config struct {
	SmtpHost       string
	SmtpPort       string
	EmailSender    string
	PasswordSender string
	NameSender     string

	Participants map[string]string   // Участники [email]имя
	Restrictions map[string][]string // Ограничения [emailSender]emailRecipient
}

var cfg Config

func init() {
	if _, err := toml.DecodeFile(configFilename, &cfg); err != nil {
		exitOnErr(err)
	}
}

func main() {
	exec()
}

func exec() {
	c, err := calculator.NewCalculator(cfg.Restrictions)
	if err != nil {
		fmt.Printf("при расчёте произошоа ошибка: %v", err)
	}

	res := c.CalculateRecipient()
	wg := sync.WaitGroup{}
	start := time.Now()

	fmt.Printf("Начинается рассылка email\n")

	for senderID, recipientID := range res {
		senderID, recipientID := senderID, recipientID //copy + shadow
		wg.Add(1)

		go func(senderID, recipientID string) {
			defer wg.Done()

			body := fmt.Sprintf("Вы участвуете в тайном Санте. Получателем Вашего подарка будет: %s", cfg.Participants[recipientID])
			recipient := mail.Address{Name: cfg.Participants[senderID], Address: senderID}
			sender := mail.Address{Name: cfg.NameSender, Address: cfg.EmailSender}
			if err := sendMail(recipient, sender, cfg.PasswordSender, cfg.SmtpHost, cfg.SmtpPort, body); err != nil {
				fmt.Printf("Произошла ошибка при отправке email на адрес: %s, ошибка: %v\n", recipient.Address, err)
			} else {
				fmt.Printf("Письмо успешно отправлено на адрес: %s\n", recipient.Address)
			}
		}(string(senderID), string(recipientID))
	}

	wg.Wait()
	fmt.Printf("Рассылка окончена. Завершена за: %s, всего email отправлено: %d\n", time.Since(start), len(res))
}

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func sendMail(to, from mail.Address, passwordSender, smtpHost, smtpPort, body string) error {
	auth := smtp.PlainAuth(
		"",
		from.Address,
		passwordSender,
		smtpHost,
	)
	title := "Тайный Санта"

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = mime.QEncoding.Encode("UTF-8", title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	// Сообщение.
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
}

// sendMailTest Вызывать для проверки/дебага вместо sendMail
//
// nolint:unused
func sendMailTest(to, from mail.Address, passwordSender, smtpHost, smtpPort, body string) error {
	fmt.Printf("Имитация отправки письма. Получатель письма: %s %s. Текст: %s\n", to.Address, to.Name, body)
	time.Sleep(1 * time.Second)
	return nil
}
