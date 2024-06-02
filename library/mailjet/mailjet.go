package mailjet

import (
	"fmt"
	"log"
	"strings"

	"luxe-beb-go/configs"
	"luxe-beb-go/library/templatehtml"

	"github.com/PuerkitoBio/goquery"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type ContentMailjet struct {
	Content  string
	To       string
	ToName   string
	Cc       string
	CcName   string
	Bcc      string
	BccName  string
	From     string
	FromName string
	Subject  string
	TextPart string
}

// sales
func SendMailInvoice(content ContentMailjet) {
	config, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalln("failed to get configuration: ", errConfig)
	}

	var contentMailjet ContentMailjet
	contentMailjet.Content = GetInvoiceTemplate()
	contentMailjet.To = content.To
	contentMailjet.ToName = content.ToName
	contentMailjet.TextPart = ""
	contentMailjet.Subject = "Subject here"

	SendMail(config, contentMailjet)
}

func SendMail(config *configs.Config, contentMailjet ContentMailjet) error {
	var err error
	var content string

	content = contentMailjet.Content
	contentMailjet.From = config.MjSenderEmail
	contentMailjet.FromName = config.MjSenderName
	mailjetClient := mailjet.NewMailjetClient(config.MjApikeyPublic, config.MjApikeyPrivate)
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: contentMailjet.From,
				Name:  contentMailjet.FromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: contentMailjet.To,
					Name:  contentMailjet.ToName,
				},
			},
			Subject:  contentMailjet.Subject,
			TextPart: contentMailjet.TextPart,
			HTMLPart: content,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return err
	}

	fmt.Printf("Data: %+v\n", res)

	return nil
}

func GetInvoiceTemplate() string {
	sampleHTML := templatehtml.TemplateInvoice()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sampleHTML))
	if err != nil {
		log.Fatal(err)
	}

	modifiedHTML, err := goquery.OuterHtml(doc.Selection)
	if err != nil {
		log.Fatal(err)
	}

	// str := html.EscapeString(modifiedHTML)

	// // for dev
	// f, err := os.Create("check.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer f.Close()

	// _, err2 := f.WriteString(modifiedHTML)
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }

	// fmt.Println("done")

	return modifiedHTML
}
