package db

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "api.db")
	if err != nil {
		panic("can't connect to DB")
	}
	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(15)
	createTables()
}
func createTables() {
	createUserTable := `
	CREATE TABLE  IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`
	_, err := DB.Exec(createUserTable)
	if err != nil {
		panic("could not create users table.")
	}
	createGroupTable := `
	CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_date DATETIME NOT NULL
	)
	`
	_, err = DB.Exec(createGroupTable)
	if err != nil {
		panic("could not create groups table.")
	}

	createGroupParticipantsTable := `
	CREATE TABLE IF NOT EXISTS group_participants (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER,
		user_id INTEGER,
		role INTEGER,
		FOREIGN KEY(group_id) REFERENCES groups(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createGroupParticipantsTable)
	if err != nil {
		panic("could not create group participants table.")
	}
	createExpenseTable := `
	CREATE TABLE IF NOT EXISTS expense (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		amount INTEGER,
		group_id INTEGER,
		payor_id INTEGER,
		created_date DATETIME NOT NULL,
		FOREIGN KEY(group_id) REFERENCES groups(id),
		FOREIGN KEY(payor_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createExpenseTable)
	if err != nil {
		panic("could not create expense table.")
	}
	createExpenseMembersTable := `
	CREATE TABLE IF NOT EXISTS expense_members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expense_id INTEGER,
		user_id INTEGER,
		owed_amount INTEGER,
		FOREIGN KEY(expense_id) REFERENCES expense(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createExpenseMembersTable)
	if err != nil {
		panic("could not create  expense_members table.")
	}
	createUserToUserOwesTable := `
	CREATE TABLE IF NOT EXISTS user_to_user_owes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER,
		payor INTEGER,
		payee INTEGER,
		debt_amount FLOAT,
		FOREIGN KEY(group_id) REFERENCES groups(id),
		FOREIGN KEY(payor) REFERENCES users(id),
		FOREIGN KEY(payee) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createUserToUserOwesTable)
	if err != nil {
		panic("could not create  user_to_user_owes table.")
	}
}
