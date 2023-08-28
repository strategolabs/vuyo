package gomail

import (
	"bytes"
	"html/template"
	"log"

	gomail "gopkg.in/gomail.v2"
)

func test() {
	// Email configuration
	from := "sender@example.com"
	to := "recipient@example.com"
	subject := "Test Email"
	smtpHost := "smtp.example.com"
	smtpPort := 587
	smtpUsername := "your_smtp_username"
	smtpPassword := "your_smtp_password"

	// Load the HTML email template
	htmlTemplate, err := template.ParseFiles("email_template.html")
	if err != nil {
		log.Fatalf("Error parsing HTML template: %v", err)
	}

	// Create a buffer to hold the generated HTML content
	var emailBody bytes.Buffer
	data := struct {
		Name string
	}{
		Name: "John Doe", // Replace this with the dynamic content you want
	}

	// Execute the template with the provided data and write the result to the buffer
	err = htmlTemplate.Execute(&emailBody, data)
	if err != nil {
		log.Fatalf("Error executing HTML template: %v", err)
	}

	// Send the email using gomail
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody.String())

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Error sending email: %v", err)
	}

	log.Println("Email sent successfully!")
}
