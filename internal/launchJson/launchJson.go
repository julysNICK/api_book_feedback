package launchjson

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ApplicationLog struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
type JsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content any    `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type BookData struct {
	ID          int    `json:"-"`
	Title       string `json:"title,omitempty"`
	NameAuthor  string `json:"name_author,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	YearRelease string `json:"year_release,omitempty"`
}

type Book_info struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	NameAuthor  string    `json:"name_author"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	YearRelease any       `json:"year_release"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (app *ApplicationLog) LaunchJsonSuccessUser(w http.ResponseWriter, r *http.Request, messageData string, contentData BookData) {
	j := JsonResponse{
		OK:      true,
		Message: messageData,
		Content: contentData,
	}

	out, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		app.ErrorLog.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *ApplicationLog) LaunchJsonErrorUser(w http.ResponseWriter, r *http.Request, messageData string, contentData string) {
	j := JsonResponse{
		OK:      false,
		Message: messageData,
		Content: contentData,
	}

	out, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		app.ErrorLog.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *ApplicationLog) LaunchJsonSuccessBook(w http.ResponseWriter, r *http.Request, messageData string, contentData Book_info) {
	j := JsonResponse{
		OK:      true,
		Message: messageData,
		Content: contentData,
	}

	out, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		app.ErrorLog.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *ApplicationLog) LaunchJsonErrorBook(w http.ResponseWriter, r *http.Request, messageData string, contentData string) {
	j := JsonResponse{
		OK:      false,
		Message: messageData,
		Content: contentData,
	}

	out, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		app.ErrorLog.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *ApplicationLog) LaunchJsonSuccessMessage(w http.ResponseWriter, r *http.Request, messageData string, contentData string) {
	j := JsonResponse{
		OK:      true,
		Message: messageData,
		Content: contentData,
	}

	out, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		app.ErrorLog.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *ApplicationLog) LaunchJsonErrorMessage(w http.ResponseWriter, r *http.Request, messageData string, contentData string, code int) {
	j := JsonResponse{
		OK:      false,
		Message: messageData,
		Content: contentData,
	}

	out, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		app.ErrorLog.Println(err)
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
