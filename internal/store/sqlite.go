package store

import (
	"database/sql"

	"github.com/yabu1121/expense-api/internal/model"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(filePath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	sqliteStore := &SQLiteStore{
		db: db,
	}

	if err := sqliteStore.createTable(); err != nil {
		db.Close()
		return nil, err
	}
	return sqliteStore, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) createTable() error {
	_, err := s.db.Exec(
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

func (s *SQLiteStore) GetAllExpenses() ([]model.Expense, error) {
	rows, err := s.db.Query(`
		select id, title, amount, category
		from expenses
		order by id asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]model.Expense, 0)

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

func (s *SQLiteStore) GetExpenseByID(id int) (*model.Expense, error) {
	row := s.db.QueryRow(`
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

func (s *SQLiteStore) CreateExpense(expense model.Expense) (*model.Expense, error) {
	result, err := s.db.Exec(`
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

func (s *SQLiteStore) UpdateExpense(expense model.Expense) (*model.Expense, error) {
	result, err := s.db.Exec(`
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
		return nil, sql.ErrNoRows
	}

	return &expense, nil
}

func (s *SQLiteStore) DeleteExpense(id int) error {
	res, err := s.db.Exec(`
		delete from expenses
		where id = ?
	`, id)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return sql.ErrNoRows
	}
	return nil
}
