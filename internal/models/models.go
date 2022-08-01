package models

import (
	"context"
	"database/sql"
	"time"
)

// DBModel is the type for database connection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Book is the type for all Book
type Book_info struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	NameAuthor  string    `json:"name_author"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	YearRelease string    `json:"year_release"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// Book is the type for all Book
type LevelsStruct struct {
	ID            int       `json:"id"`
	Level         int       `json:"level"`
	NumberOpinion int       `json:"number_opinion"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// User is the type for all User
type User struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Surname       string    `json:"surname"`
	Password      string    `json:"password"`
	NumberOpinion int       `json:"number_Opinion"`
	LevelNumber   int       `json:"level_number"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// Opinions is the type for all Opinions
type Opinions struct {
	ID        int       `json:"id"`
	Feedback  string    `json:"feedback"`
	IDUser    int       `json:"id_user"`
	IDBook    int       `json:"id_book"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) ExistBook(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var book Book_info
	row := m.DB.QueryRowContext(ctx, `
	select 
			id, title ,name_author ,description, year_release, coalesce(image, ''),
			created_at, updated_at
		from 
			book_info 
		where id = ?`, id)

	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.NameAuthor,
		&book.Description,
		&book.YearRelease,
		&book.Image,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (m *DBModel) GetBook(id int) (Book_info, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var book Book_info
	row := m.DB.QueryRowContext(ctx, `
		select 
			id, title ,name_author ,description, year_release, coalesce(image, ''),
			created_at, updated_at
		from 
			book_info 
		where id = ?`, id)

	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.NameAuthor,
		&book.Description,
		&book.YearRelease,
		&book.Image,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return book, err
	}

	return book, nil

}

func (m *DBModel) InsertBook(c Book_info) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	insert into book_info 
					(title,  name_author, description, image, year_release,created_at, updated_at)
					values (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		c.Title,
		c.NameAuthor,
		c.Description,
		c.Image,
		c.YearRelease,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *DBModel) GetAllBooks() ([]Book_info, error) {

	var book []Book_info
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	row, err := m.DB.QueryContext(ctx, "SELECT * FROM book_info")

	if err != nil {
		println(err.Error())
		return book, err
	}

	for row.Next() {
		var id int
		var title string
		var nameAuthor string
		var description string
		var image string
		var yearRelease string
		var createdAt time.Time
		var updatedAt time.Time
		err = row.Scan(&id, &title, &nameAuthor, &description, &image, &yearRelease, &createdAt, &updatedAt)
		if err != nil {
			println(err.Error())
		}

		book = append(book, Book_info{ID: id, Title: title, NameAuthor: nameAuthor, Description: description, Image: image, YearRelease: yearRelease, CreatedAt: createdAt, UpdatedAt: updatedAt})
	}

	return book, nil
}

func (m *DBModel) DeleteBook(id int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	delete from book_info where id = ?
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		id,
	)
	if err != nil {
		return 0, nil
	}

	_, errDeleteOpinion := m.DeleteOpinion(id, "book")

	if errDeleteOpinion != nil {
		return 0, nil
	}

	return 1, nil

}

func (m *DBModel) DeleteOpinion(id int, tableNameDelete string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	switch tableNameDelete {
	case "user":
		stmt := `
	 delete from opinions where id_user = ?
	 `
		_, err := m.DB.ExecContext(ctx, stmt,
			id,
		)

		if err != nil {
			return 0, nil
		}

		return 1, nil

	case "book":
		stmt := `
	 delete from opinions where id_book = ?
	 `
		_, err := m.DB.ExecContext(ctx, stmt,
			id,
		)

		if err != nil {
			return 0, nil
		}

		return 1, nil

	case "opinion":
		stmt := `
	 delete from opinions where id = ?
	 `
		_, err := m.DB.ExecContext(ctx, stmt,
			id,
		)

		if err != nil {
			return 0, nil
		}

		return 1, nil
	default:
		return 2, nil
	}
}

//----------------------------user---------------------------------

func (m *DBModel) InsertUser(userStruct User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `	insert into users 
						(name,  email, surname, password, created_at, updated_at)
						values (?, ?, ?, ?, ?, ?)
					`

	result, err := m.DB.ExecContext(ctx, stmt,
		userStruct.Name,
		userStruct.Email,
		userStruct.Surname,
		userStruct.Password,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		println(err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (m *DBModel) FindOneUser(emailUser string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userSearchBD User
	row := m.DB.QueryRowContext(ctx, `
						SELECT * FROM users WHERE email = ?
					`, emailUser)
	err := row.Scan(
		&userSearchBD.ID,
		&userSearchBD.Name,
		&userSearchBD.Email,
		&userSearchBD.Surname,
		&userSearchBD.Password,
		&userSearchBD.NumberOpinion,
		&userSearchBD.LevelNumber,
		&userSearchBD.CreatedAt,
		&userSearchBD.UpdatedAt,
	)
	if err != nil {
		return userSearchBD, err
	}
	return userSearchBD, nil
}
func (m *DBModel) FindOneUserById(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userSearchBD User
	row := m.DB.QueryRowContext(ctx, `
						SELECT * FROM users WHERE id = ?
					`, id)
	err := row.Scan(
		&userSearchBD.ID,
		&userSearchBD.Name,
		&userSearchBD.Email,
		&userSearchBD.Surname,
		&userSearchBD.Password,
		&userSearchBD.NumberOpinion,
		&userSearchBD.LevelNumber,
		&userSearchBD.CreatedAt,
		&userSearchBD.UpdatedAt,
	)
	if err != nil {
		return userSearchBD, err
	}
	return userSearchBD, nil
}

func (m *DBModel) ExistUser(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User
	row := m.DB.QueryRowContext(ctx, `
		select 
			id, name ,email ,surname, created_at, updated_at
		from 
			users 
		where id = ?`, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Surname,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (m *DBModel) UpdateOpinionNumbers(id int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.FindOneUserById(id)
	if err != nil {
		println(err.Error())
		return 0, err
	}

	sumOpinionNumber := result.NumberOpinion + 1

	stmt := `UPDATE 
			users
		SET 
			number_opinion = ? 
		where id = ?
					`

	_, err = m.DB.ExecContext(ctx, stmt,
		sumOpinionNumber, id)

	if err != nil {
		println(err.Error())
		return 0, err
	}

	return int(id), nil

}

func (m *DBModel) FindManyLevels() ([]LevelsStruct, error) {
	var levels []LevelsStruct
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	row, err := m.DB.QueryContext(ctx, "SELECT * FROM levels")

	if err != nil {
		println(err.Error())
		return nil, err
	}

	for row.Next() {
		var id int
		var level int
		var numberOpinion int
		var createdAt time.Time
		var updatedAt time.Time
		err = row.Scan(&id, &level, &numberOpinion, &createdAt, &updatedAt)
		if err != nil {
			println(err.Error())
		}

		levels = append(levels, LevelsStruct{ID: id, Level: level, NumberOpinion: numberOpinion, CreatedAt: createdAt, UpdatedAt: updatedAt})
	}

	return levels, nil
}

func IsReadyLevelUp(resultsArrayLevel []LevelsStruct, resultUser User) []LevelsStruct {
	var levelsReady []LevelsStruct
	for result := range resultsArrayLevel {
		if resultUser.NumberOpinion == resultsArrayLevel[result].NumberOpinion {
			levelsReady = append(levelsReady, resultsArrayLevel[result])
		}
	}
	return levelsReady
}

func (m *DBModel) UpdateLevelUp(id int) (int, error) {
	resultsArrayLevel, err := m.FindManyLevels()
	if err != nil {
		println(err.Error())
		return 0, err
	}
	resultUser, err := m.FindOneUserById(id)
	if err != nil {
		println(err.Error())
		return 0, err
	}

	levelsUp := IsReadyLevelUp(resultsArrayLevel, resultUser)

	if len(levelsUp) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err != nil {
			println(err.Error())
			return 0, err
		}

		stmt := `UPDATE 
			users
		SET 
			level_number = ? 
		where id = ?
					`

		_, err = m.DB.ExecContext(ctx, stmt,
			levelsUp[0].Level, id)

		if err != nil {
			println(err.Error())
			return 0, err
		}
		return 1, nil

	} else {
		return 0, nil

	}
}

func (m *DBModel) ExistOpinion(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var opinion Opinions
	row := m.DB.QueryRowContext(ctx, `
		select 
			id, feedback ,id_user ,id_book, created_at, updated_at
		from 
			users 
		where id = ?`, id)

	err := row.Scan(
		&opinion.ID,
		&opinion.Feedback,
		&opinion.IDUser,
		&opinion.IDBook,
		&opinion.CreatedAt,
		&opinion.UpdatedAt,
	)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (m *DBModel) GetOpinionByIdBook(id int) (Opinions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var opinion Opinions
	row := m.DB.QueryRowContext(ctx, `
		select 
			id, feedback ,id_user ,id_book, created_at, updated_at
		from 
			opinions 
		where id_book = ?`, id)

	err := row.Scan(
		&opinion.ID,
		&opinion.Feedback,
		&opinion.IDUser,
		&opinion.IDBook,
		&opinion.CreatedAt,
		&opinion.UpdatedAt,
	)
	if err != nil {
		return opinion, err
	}

	return opinion, nil

}
func (m *DBModel) GetOpinionByIdUser(id int) (Opinions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var opinion Opinions
	row := m.DB.QueryRowContext(ctx, `
		select 
			id, feedback ,id_user ,id_book, created_at, updated_at
		from 
			opinions 
		where id_user = ?`, id)

	err := row.Scan(
		&opinion.ID,
		&opinion.Feedback,
		&opinion.IDUser,
		&opinion.IDBook,
		&opinion.CreatedAt,
		&opinion.UpdatedAt,
	)
	if err != nil {
		return opinion, err
	}

	return opinion, nil

}
func (m *DBModel) InsertOpinions(opinion Opinions) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `	insert into opinions 
						(feedback,  id_user, id_book, created_at, updated_at)
						values (?, ?, ?, ?, ?)
					`

	result, err := m.DB.ExecContext(ctx, stmt,
		opinion.Feedback,
		opinion.IDUser,
		opinion.IDBook,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		println(err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (m *DBModel) GetAllOpinions() ([]Opinions, error) {
	var opinion []Opinions
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()
	row, err := m.DB.QueryContext(ctx, "SELECT * FROM opinions")

	if err != nil {
		println(err.Error())
		return opinion, err
	}

	for row.Next() {
		var id int
		var feedback string
		var id_user int
		var id_book int

		var createdAt time.Time
		var updatedAt time.Time
		err = row.Scan(&id, &feedback, &id_user, &id_book, &createdAt, &updatedAt)
		if err != nil {
			println(err.Error())
		}

		opinion = append(opinion, Opinions{ID: id, Feedback: feedback, IDUser: id_user, IDBook: id_book, CreatedAt: createdAt, UpdatedAt: updatedAt})
	}

	return opinion, nil
}
