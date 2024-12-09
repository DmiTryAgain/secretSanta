//nolint:dupl
package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/mail"
	"net/smtp"
	"os"
	"sync"
	"time"

	"github.com/DmiTryAgain/secretSanta/pkg/calculator"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

type App struct {
	rootCmd *cobra.Command
	cfg     Config
}

func New(rootCmd *cobra.Command) App {
	return App{
		rootCmd: rootCmd,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.rootCmd.AddCommand(a.createConfigCmd(), a.runCmd(), a.dryRunCmd())

	return a.rootCmd.Execute()
}

func (a *App) createConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "createConfig",
		Short: "Create a new config file",
		Long: `The command create a config file "conf.toml" which is used by "run" command.
Config file has some sections to set up your configuration.
Firstly, you must set your smtp server data to send results data.
Also, you must set login, password and sender name from your smtp server configuration.
Than, you must set "Participants" section. It looks like mapping email->name of participants.
Email is using to send him result of calculation.
Name is using in email body to other recipient who is getting current participant as a result of calculation.
Emails must be unique.
And the last section is mapping email->[email] of restrictions. Key - is participant email which we set restrictions.
Values - list of other emails which will never have been matched to current participant.
You shouldn't set restriction for each email itself, because it will have already been during calculation.
	`,
		Run: func(cmd *cobra.Command, args []string) {
			var buf bytes.Buffer
			// write default DB config
			buf.WriteString(`## smtp сервер конфигурация
SmtpHost = "smtp.yandex.ru"
SmtpPort = "587"

## Информация об отправителе
EmailSender = "example@yandex.ru"
NameSender = "exampleName"
PasswordSender = "examplePass"

## Информация об участниках email->имя
[Participants]
"email1@test.com" = "Петя"
"email2@test.com" = "Света"
"email3@test.com" = "Вася"
"email4@test.com" = "Аня"
"email5@test.com" = "Коля"
"email6@test.com" = "Света"
"email7@test.com" = "Борис"

## Ограничения email отправителя, которому не должен достаться email получателя emailSender->emailRecipient
[Restrictions]
"email1@test.com" = [ "email2@test.com" ]
"email2@test.com" = [ "email1@test.com" ]
"email3@test.com" = [ "email4@test.com" ]
"email4@test.com" = [ "email3@test.com" ]
"email5@test.com" = [ "email6@test.com" ]
"email6@test.com" = [ "email5@test.com", "email1@test.com" ]
"email7@test.com" = []`)
			if err := os.WriteFile(configFilename, buf.Bytes(), 0o644); err != nil {
				log.Fatalf("Failed to write file %s: %v", configFilename, err)
				return
			}

			log.Printf("File %v was successfully created.", configFilename)
		},
	}
}

func (a *App) runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run calculator with real results sending",
		Long: `The command reads config, initializes calculator, validates participants and restricts.
Than it calculates results and runs real result sending.
It prints log to console and if something is wrong you can see it`,
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := toml.DecodeFile(configFilename, &a.cfg); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					log.Fatal("configuration file was not found. Please create new via `secretSanta createConfig`")
				}

				log.Fatal(errors.Wrap(err, "decode config file failed"))
			}

			c, err := calculator.NewCalculator(a.cfg.Participants, a.cfg.Restrictions)
			if err != nil {
				log.Fatal(errors.Wrap(err, "init calculator failed"))
			}

			res := c.CalculateRecipient()
			wg := sync.WaitGroup{}
			start := time.Now()

			fmt.Printf("Начинается рассылка email\n")
			for senderID, recipientID := range res {
				senderID, recipientID := senderID, recipientID // copy + shadow
				wg.Add(1)

				go func(senderID, recipientID string) {
					defer wg.Done()
					body := fmt.Sprintf("Вы участвуете в тайном Санте. Получателем Вашего подарка будет: %s", a.cfg.Participants[recipientID])
					participantAddress := mail.Address{Name: a.cfg.Participants[senderID], Address: senderID}
					if err := a.sendMail(participantAddress, body); err != nil {
						fmt.Printf("Произошла ошибка при отправке email на адрес: %s, ошибка: %v\n", participantAddress.Address, err)
					} else {
						fmt.Printf("Письмо успешно отправлено на адрес: %s\n", participantAddress.Address)
					}
				}(string(senderID), string(recipientID))
			}

			wg.Wait()
			fmt.Printf("Рассылка окончена. Завершена за: %s, всего email отправлено: %d\n", time.Since(start), len(res))
		},
	}
}

func (a *App) dryRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dryRun",
		Short: "Run calculator with result sending imitation",
		Long: `The command reads config, initializes calculator, validates participants and restricts.
Than it calculates results and runs imitation result sending with logs.
It may be useful to run for validating config and calculating before real sending by "run" command.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Читаем конфиг
			if _, err := toml.DecodeFile(configFilename, &a.cfg); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					log.Fatal("configuration file was not found. Please create new via `secretSanta createConfig`")
				}

				log.Fatal(errors.Wrap(err, "decode config file failed"))
			}

			// Инициируем калькулятор
			c, err := calculator.NewCalculator(a.cfg.Participants, a.cfg.Restrictions)
			if err != nil {
				log.Fatal(errors.Wrap(err, "init calculator failed"))
			}

			// Рассчитываем
			res := c.CalculateRecipient()

			// Отправляем результаты
			wg := sync.WaitGroup{}
			start := time.Now()
			fmt.Printf("Начинается рассылка email\n")
			for senderID, recipientID := range res {
				senderID, recipientID := senderID, recipientID // copy + shadow
				wg.Add(1)

				go func(senderID, recipientID string) {
					defer wg.Done()
					body := fmt.Sprintf("Вы участвуете в тайном Санте. Получателем Вашего подарка будет: %s", a.cfg.Participants[recipientID])
					participantAddress := mail.Address{Name: a.cfg.Participants[senderID], Address: senderID}
					if err := a.sendMailTest(participantAddress, body); err != nil {
						fmt.Printf("Произошла ошибка при отправке email на адрес: %s, ошибка: %v\n", participantAddress.Address, err)
					} else {
						fmt.Printf("Письмо успешно отправлено на адрес: %s\n", participantAddress.Address)
					}
				}(string(senderID), string(recipientID))
			}

			wg.Wait()
			fmt.Printf("Рассылка окончена. Завершена за: %s, всего email отправлено: %d\n", time.Since(start), len(res))
		},
	}
}

// sendMail Отправляет письмо получателю participantAddress с телом body
func (a *App) sendMail(participantAddress mail.Address, body string) error {
	// Отправитель будет всегда один и тот же, берётся из конфига
	senderAddress := mail.Address{Name: a.cfg.NameSender, Address: a.cfg.EmailSender}

	// Инициируем структуру для аутентификации
	auth := smtp.PlainAuth(
		"",
		senderAddress.Address,
		a.cfg.PasswordSender,
		a.cfg.SmtpHost,
	)
	title := "Тайный Санта"

	// Заголовки, чтобы почтовый сервис не счёл нас спам-рассылкой
	header := make(map[string]string)
	header["From"] = senderAddress.String()
	header["To"] = participantAddress.String()
	header["Subject"] = mime.QEncoding.Encode("UTF-8", title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	// Сообщение, наполняется сначала заголовками, потом - контентом
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	return smtp.SendMail(
		a.cfg.SmtpHost+":"+a.cfg.SmtpPort,
		auth,
		senderAddress.Address,
		[]string{participantAddress.Address},
		[]byte(message),
	)
}

// sendMailTest Вызывается для проверки/дебага вместо sendMail
func (a *App) sendMailTest(participantAddress mail.Address, body string) error {
	fmt.Printf("Имитация отправки письма. Получатель письма: %s %s. Текст: %s\n", participantAddress.Address, participantAddress.Name, body)
	time.Sleep(1 * time.Second)
	return nil
}
