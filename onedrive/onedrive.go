package onedrive

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hiteshrepo/ms-graph-api/token"
)

func UploadFile() {
	accessToken, err := token.GetTokenByPublicClientApp()
	if err != nil {
		panic(err)
	}

	uploadPath := "/Documents/"

	fileName := "example.txt"

	fileContent := []byte("This is the content of the uploaded file.")

	uploadURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/root:%s%s:/content", uploadPath, fileName)

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileContent))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error uploading file:", resp.Status)
		return
	}

	fmt.Println("File uploaded:", fileName)
}

func DownloadFile(downloadURL string) {
	resp, err := http.Get(downloadURL)
	if err != nil {
		fmt.Println("Error making download request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error downloading file:", resp.Status)
		return
	}

	localFilePath := "downloaded_file.txt"

	file, err := os.Create(localFilePath)
	if err != nil {
		fmt.Println("Error creating local file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error saving file content to local file:", err)
		return
	}

	fmt.Println("File downloaded and saved to:", localFilePath)

}
