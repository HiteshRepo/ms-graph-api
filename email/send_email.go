package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hiteshrepo/ms-graph-api/token"
)

type Email struct {
	Message `json:"message"`
}

type Message struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	Body    struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	} `json:"body"`
	ToRecipients []struct {
		EmailAddress struct {
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"toRecipients"`
	Attachments []struct {
		ODataType string `json:"@odata.type"`
		Name      string `json:"name"`
		Content   string `json:"contentBytes"`
	} `json:"attachments"`
}

func SendEmail() {
	email := Email{}
	email.Message.Subject = "Test Email with Attachment"
	email.Message.Body.ContentType = "HTML"
	email.Message.Body.Content = "<html><body>This is the email body.</body></html>"
	email.Message.ToRecipients = []struct {
		EmailAddress struct {
			Address string `json:"address"`
		} `json:"emailAddress"`
	}{
		{
			EmailAddress: struct {
				Address string `json:"address"`
			}{
				Address: "recipient@example.com",
			},
		},
	}

	email.Message.Attachments = []struct {
		ODataType string `json:"@odata.type"`
		Name      string `json:"name"`
		Content   string `json:"contentBytes"`
	}{
		{
			ODataType: "#microsoft.graph.fileAttachment",
			Name:      "example.txt",
			Content:   "base64-encoded-attachment-content",
		},
	}

	emailJSON, err := json.Marshal(email)
	if err != nil {
		panic(err)
	}

	token, err := token.GetTokenByPublicClientApp()
	if err != nil {
		panic(err)
	}

	sendEmailUsingGraphAPI(token, emailJSON)
}

func FetchEmailMessages() []Message {
	accessToken, err := token.GetTokenByPublicClientApp()
	if err != nil {
		panic(err)
	}

	messagesURL := "https://graph.microsoft.com/v1.0/me/messages"

	req, err := http.NewRequest("GET", messagesURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error retrieving messages:", resp.Status)
		return nil
	}

	var response struct {
		Messages []Message `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding message response:", err)
		return nil
	}

	for _, msg := range response.Messages {
		fmt.Println("Subject:", msg.Subject)
	}

	return response.Messages

}

func FetchAnEmailMessage() {

	allMessages := FetchEmailMessages()
	if len(allMessages) == 0 {
		panic("no messages")
	}

	aMessage := allMessages[0]

	messageID := aMessage.ID

	accessToken, err := token.GetTokenByPublicClientApp()
	if err != nil {
		panic(err)
	}

	emailURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s", messageID)
	req, err := http.NewRequest("GET", emailURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	var emailData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emailData); err != nil {
		fmt.Println("Error decoding email response:", err)
		return
	}

	attachments := emailData["attachments"].([]interface{})
	for _, attachment := range attachments {
		attachmentInfo := attachment.(map[string]interface{})
		attachmentName := attachmentInfo["name"].(string)
		attachmentID := attachmentInfo["id"].(string)

		downloadAttachment(attachmentID, accessToken, attachmentName, messageID)
	}
}

func downloadAttachment(attachmentID, accessToken, attachmentName, messageID string) {
	attachmentURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s/attachments/%s/$value", messageID, attachmentID)
	req, err := http.NewRequest("GET", attachmentURL, nil)
	if err != nil {
		fmt.Println("Error creating request for attachment:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error downloading attachment:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error downloading attachment:", resp.Status)
		return
	}

	file, err := os.Create(attachmentName)
	if err != nil {
		fmt.Println("Error creating attachment file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error saving attachment content to file:", err)
		return
	}

	fmt.Println("Attachment downloaded:", attachmentName)
}

func sendEmailUsingGraphAPI(accessToken string, emailJSON []byte) error {
	req, err := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/me/sendMail", bytes.NewBuffer(emailJSON))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body := []byte{}
	_, err = resp.Body.Read(body)
	if err != nil {
		return err
	}

	fmt.Println("Response from api: ", string(body))

	return nil
}
