package main

import (
	"encoding/json"
	"log"
	"mybook/internal/authentication"
	launchjson "mybook/internal/launchJson"
	"mybook/internal/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	ID       int    `json:"-"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Password string `json:"password,omitempty"`
}
type OpinionData struct {
	ID       int    `json:"-"`
	Feedback string `json:"feedback"`
	IDUser   int    `json:"id_user"`
	IDBook   int    `json:"id_book"`
}
type UserJson struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Password string `json:"-"`
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

type Token struct {
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content any    `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type jsonResponseWithJwt struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content any    `json:"content,omitempty"`
	Token   any    `json:"token,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		jsonFormat := launchjson.ApplicationLog{
			InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
			ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		if r.Header["Authorization"] == nil {

			jsonFormat.LaunchJsonErrorMessage(w, r, "Jwt Não encontrado", "", 401)
			return
		}

		tokenValid, err := authentication.ValidJwt(r.Header["Authorization"][0])

		if err != nil {
			jsonFormat.LaunchJsonErrorMessage(w, r, "Seu token expirou", err.Error(), 401)
			return

		}

		if tokenValid == "Jwt válido" {
			r.Header.Set("authenticated", "true")
			handler.ServeHTTP(w, r)
			return

		}

		jsonFormat.LaunchJsonErrorMessage(w, r, "Não autorizado", "Revise suas credenciais", 401)

	}
}

func (app *application) GetInitial(w http.ResponseWriter, r *http.Request) {
	out, err := json.MarshalIndent("welcome to the book reviews api", "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

func (app *application) GetBookById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jsonFormat := launchjson.ApplicationLog{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	bookId, _ := strconv.Atoi(id)

	book, err := app.DB.GetBook(bookId)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			jsonFormat.LaunchJsonErrorMessage(w, r, "Não existe livro com esse id", err.Error(), 400)
			return
		}

		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro", err.Error(), 400)
		app.errorLog.Println(err.Error())
		return
	}
	out, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
func (app *application) SaveBookHandle(w http.ResponseWriter, r *http.Request) {
	var txnData BookData
	jsonFormat := launchjson.ApplicationLog{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	err := json.NewDecoder(r.Body).Decode(&txnData)
	if err != nil {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro", err.Error(), 400)
		app.errorLog.Println(err)
		return
	}
	txnData = BookData{
		Title:       txnData.Title,
		NameAuthor:  txnData.NameAuthor,
		Description: txnData.Description,
		Image:       txnData.Image,
		YearRelease: txnData.YearRelease,
	}

	app.SaveBook(txnData.Title, txnData.NameAuthor, txnData.Description, txnData.Image, txnData.YearRelease)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	jsonFormat.LaunchJsonSuccessUser(w, r, "Livro salvo com sucesso", launchjson.BookData(txnData))

}

func (app *application) ListBooksHandle(w http.ResponseWriter, r *http.Request) {

	txnData, err := app.DB.GetAllBooks()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	j := jsonResponse{
		OK:      true,
		Message: "Lista de livros",
		Content: txnData,
	}

	out, err := json.MarshalIndent(j, "", "   ")

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) SaveBook(title, nameAuthor, description, image, year_release string) (int, error) {

	book := models.Book_info{
		Title:       title,
		NameAuthor:  nameAuthor,
		Description: description,
		Image:       image,
		YearRelease: year_release,
	}

	id, err := app.DB.InsertBook(book)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (app *application) SaveUserHandle(w http.ResponseWriter, r *http.Request) {
	var txnData UserData
	jsonFormat := launchjson.ApplicationLog{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	err := json.NewDecoder(r.Body).Decode(&txnData)
	if err != nil {
		app.errorLog.Println(err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(txnData.Password), bcrypt.DefaultCost)

	txnData = UserData{
		Name:     txnData.Name,
		Email:    txnData.Email,
		Surname:  txnData.Surname,
		Password: string(hash),
	}
	id, errUser := app.SaveUser(txnData.Name, txnData.Email, txnData.Surname, txnData.Password)

	if errUser != nil {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Verifique seu email ou senha", "", 400)
		app.errorLog.Println(err)
		return
	}

	jsonFormat.LaunchJsonSuccessMessage(w, r, "Usuario cadastrado com sucesso", strconv.Itoa(id))

}

func (app *application) SaveUser(name, email, surname, password string) (int, error) {
	user := models.User{
		Name:     name,
		Email:    email,
		Surname:  surname,
		Password: password,
	}

	id, err := app.DB.InsertUser(user)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (app *application) LoginUserHandle(w http.ResponseWriter, r *http.Request) {
	var txnData UserData
	jsonFormat := launchjson.ApplicationLog{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	err := json.NewDecoder(r.Body).Decode(&txnData)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	result, err := app.DB.FindOneUser(txnData.Email)

	if err != nil {

		jsonFormat.LaunchJsonErrorMessage(w, r, "Email ou senha errada", err.Error(), 400)
		app.errorLog.Println(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(txnData.Password)); err != nil {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Senha errada", err.Error(), 400)
		app.errorLog.Println(err)
		return
	}

	jwtToken, err := authentication.GenerateJwt(txnData.Email, result.ID)

	if err != nil {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Erro na geração do JWT", err.Error(), 400)
		app.errorLog.Println(err)
		return
	}

	j := jsonResponseWithJwt{
		OK:      true,
		Message: "Login feito com sucesso",
		Content: UserJson{
			ID:      result.ID,
			Name:    result.Name,
			Surname: result.Surname,
		},
		Token: jwtToken,
	}

	out, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		jsonFormat.LaunchJsonErrorMessage(w, r, "Erro na geração do Json", err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// opinions
func (app *application) GetOpinionByIdBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jsonFormat := launchjson.ApplicationLog{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	bookId, _ := strconv.Atoi(id)

	opinion, err := app.DB.GetOpinionByIdBook(bookId)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			jsonFormat.LaunchJsonErrorMessage(w, r, "Não existe opinião com esse id", err.Error(), 400)
			return
		}

		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro", err.Error(), 400)
		app.errorLog.Println(err.Error())
		return
	}
	out, err := json.MarshalIndent(opinion, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) GetOpinionByIdUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jsonFormat := launchjson.ApplicationLog{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	userId, _ := strconv.Atoi(id)

	opinion, err := app.DB.GetOpinionByIdUser(userId)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			jsonFormat.LaunchJsonErrorMessage(w, r, "Não existe opinião com esse id", err.Error(), 400)
			return
		}

		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro", err.Error(), 400)
		app.errorLog.Println(err.Error())
		return
	}
	out, err := json.MarshalIndent(opinion, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) DeleteBookId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jsonFormat := launchjson.ApplicationLog{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	userId, _ := strconv.Atoi(id)

	_, err := app.DB.DeleteBook(userId)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			jsonFormat.LaunchJsonErrorMessage(w, r, "Não existe livro com esse id", err.Error(), 400)
			return
		}

		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro", err.Error(), 400)
		app.errorLog.Println(err.Error())
		return
	}
	out, err := json.MarshalIndent("Livro exlcuido com sucesso", "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) ListOpinionsHandle(w http.ResponseWriter, r *http.Request) {

	txnData, err := app.DB.GetAllOpinions()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	j := jsonResponse{
		OK:      true,
		Message: "Lista de opiniões",
		Content: txnData,
	}

	out, err := json.MarshalIndent(j, "", "   ")

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) SaveOpinionHandle(w http.ResponseWriter, r *http.Request) {
	var txnData OpinionData
	jsonFormat := launchjson.ApplicationLog{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	err := json.NewDecoder(r.Body).Decode(&txnData)
	if err != nil {
		app.errorLog.Println(err)
	}

	if len(txnData.Feedback) == 0 && strconv.Itoa(txnData.IDBook) == "0" && strconv.Itoa(txnData.IDUser) == "0" {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Todos os campos são obrigatorios", "", 400)
		return
	}
	if len(txnData.Feedback) == 0 {
		jsonFormat.LaunchJsonErrorMessage(w, r, "O feedback tem que está preenchido", "", 400)
		return
	}

	if strconv.Itoa(txnData.IDBook) == "0" {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Insira o id do livro", "", 400)
		return
	}
	if strconv.Itoa(txnData.IDUser) == "0" {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Insira o id do usuario", "", 400)
		return
	}

	_, errBook := app.DB.ExistBook(txnData.IDBook)
	if errBook != nil {
		if errBook.Error() == "sql: no rows in result set" {
			jsonFormat.LaunchJsonErrorMessage(w, r, "id do livro inexistente", "", 400)
			app.errorLog.Println(errBook.Error())
			return
		}
		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro no registro do feedback", "", 400)
		app.errorLog.Println(errBook.Error())
		return
	}
	_, errUserExist := app.DB.ExistUser(txnData.IDUser)

	if errUserExist != nil {
		if errUserExist.Error() == "sql: no rows in result set" {
			jsonFormat.LaunchJsonErrorMessage(w, r, "id do usuario inexistente", "", 400)
			app.errorLog.Println(errUserExist.Error())
			return
		}
		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro no registro do feedback", "", 400)
		app.errorLog.Println(errUserExist.Error())

		return
	}
	txnData = OpinionData{
		Feedback: txnData.Feedback,
		IDUser:   txnData.IDUser,
		IDBook:   txnData.IDBook,
	}
	id, errUser := app.SaveOpinion(txnData.Feedback, strconv.Itoa(txnData.IDUser), strconv.Itoa(txnData.IDBook))
	if errUser != nil {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro no registro do feedback", "", 400)
		app.errorLog.Println(err)
		return
	}
	_, _ = app.DB.UpdateOpinionNumbers(txnData.IDUser)
	_, _ = app.DB.UpdateLevelUp(txnData.IDUser)
	if errUser != nil {
		jsonFormat.LaunchJsonErrorMessage(w, r, "Ocorreu um erro no registro do feedback", "", 400)
		app.errorLog.Println(err)
		return
	}
	jsonFormat.LaunchJsonSuccessMessage(w, r, "Opinião cadastrada com sucesso", strconv.Itoa(id))
}

func (app *application) SaveOpinion(feedback, id_user, id_book string) (int, error) {

	convIdUser, err := strconv.Atoi(id_user)

	if err != nil {
		return 0, err
	}

	convIdIdBook, err := strconv.Atoi(id_book)

	if err != nil {
		return 0, err
	}

	opinion := models.Opinions{
		Feedback: feedback,
		IDUser:   convIdUser,
		IDBook:   convIdIdBook,
	}

	id, err := app.DB.InsertOpinions(opinion)

	if err != nil {
		return 0, err
	}

	return id, nil

}
