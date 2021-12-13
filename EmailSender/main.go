package main

import (
	"context"
	"emailsender/config"
	"emailsender/db"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"sync"
	"time"
)

func main() {
	cfg := config.Init()
	dbService, err := db.NewPgxCon(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(time.Second * 15)
	for {
		select {
		case <-ticker.C:
			{
				err := StartMessaging(dbService)
				if err != nil {
					fmt.Printf("ошибка отправки сообщения %v", err)
					fmt.Println()
				}
			}
		}
	}
}

func StartMessaging(db db.DBService) error {
	list, err := db.GetReminds()
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	for _, value := range list {
		wg.Add(1)
		sendValue := value
		go func() {
			err := sendEmail(sendValue, wg)
			if err != nil {
				log.Println("ошибка при отправке напоминания на почту")
			}
			err = db.UpdateStatusReminds(sendValue.Id)
			if err != nil {
				log.Println("не удалось обновить статус напоминания на почту")
			}
		}()
	}
	wg.Wait()
	return nil
}

func sendEmail(rem db.Remind, wg *sync.WaitGroup) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "mvp.internet.technology@mail.ru")
	m.SetHeader("To", rem.Where)
	m.SetHeader("Subject", "Remind from our service")
	m.SetBody("text/html", rem.Message)

	d := gomail.NewDialer("smtp.mail.ru", 465, "mvp.internet.technology", "HzDA0idjBAjnFSwGED7T")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		wg.Done()
		return err
	}
	wg.Done()
	return nil
}
