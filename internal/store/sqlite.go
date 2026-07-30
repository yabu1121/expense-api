package store

import (
	"database/sql"
	"errors"

	"github.com/yabu1121/expense-api/internal/model"
	"golang.org/x/tools/go/analysis/passes/defers"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	DB = db

	if err := CreateTable(); err != nil {
		return err
	}
	return nil
}

func CreateTable() error {
	_, err := DB.Exec(
		`CREATE TABLE IF NOT EXISTS expenses(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			amount INTEGER,
			category TEXT
		)`,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllExpenses() ([]model.Expense, error) {
	rows, err := DB.Query(`
		select id, title, amount, category
		from expenses
		order by id asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []model.Expense

	for rows.Next() {
		var expense model.Expense
		err := rows.Scan(
			&expense.ID,
			&expense.Title,
			&expense.Amount,
			&expense.Category,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

func GetExpenseByID(id int) (*model.Expense, error) {
	row := DB.QueryRow(`
		select id, title, amount, category
		from expenses
		where id = ?
	`, id)
	var expense model.Expense
	err := row.Scan(
		&expense.ID,
		&expense.Title,
		&expense.Amount,
		&expense.Category,
	)
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func CreateExpense(expense model.Expense) (*model.Expense, error) {
	result, err := DB.Exec(`
		insert into expenses (title, amount, category)
		values(?, ?, ?)
	`, expense.Title, expense.Amount, expense.Category)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	expense.ID = int(id)

	return &expense, nil
}

func UpdateExpense(expense model.Expense) (*model.Expense, error) {
	result, err := DB.Exec(`
		update expenses
		set title = ?, amount = ?, category = ?
		where id = ?
	`, expense.Title, expense.Amount, expense.Category, expense.ID)
	if err != nil {
		return nil, err
	}

	num, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, errors.New("expense not found")
	}

	return &expense, nil
}

func DeleteExpense(id int) error {
	res, err := DB.Exec(`
		delete from expenses
		where id = ?
	`, id)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected();
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("expense not found")
	}
	return nil
}
