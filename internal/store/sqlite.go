package store

import (
	"database/sql"
	"errors"
	"uuid"

	"github.com/yabu1121/expense-api/internal/model"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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

	if err := sqliteStore.createCategoryTable(); err != nil {
		db.Close()
		return nil, err
	}

	if err := sqliteStore.createExpenseTable(); err != nil {
		db.Close()
		return nil, err
	}
	return sqliteStore, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) createCategoryTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) createExpenseTable() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS expenses (
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

func (s *SQLiteStore) ListExpenses(filter model.ExpenseFilter) ([]model.Expense, error) {
	query := `
		select id, title, amount, category
		from expenses
	`
	args := make([]any, 0, 3)

	if filter.Category != "" {
		query += " where category = ?"
		args = append(args, filter.Category)
	}

	query += " order by id asc"

	if filter.Limit > 0 {
		query += " limit ?"
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			query += " offset ?"
			args = append(args, filter.Offset)
		}
	}

	rows, err := s.db.Query(query, args...)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrExpenseNotFound
		}
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
		return nil, model.ErrExpenseNotFound
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
		return model.ErrExpenseNotFound
	}
	return nil
}

func (s *SQLiteStore) GetExpenseSummary() (*model.ExpenseSummary, error) {
	row := s.db.QueryRow(`
		select count(*), coalesce(sum(amount), 0)
		from expenses
	`)

	var summary model.ExpenseSummary
	if err := row.Scan(
		&summary.Count,
		&summary.TotalAmount,
	); err != nil {
		return nil, err
	}
	return &summary, nil
}

func (s *SQLiteStore) CreateCategory(category model.Category) (*model.Category, error) {
	category.ID = uuid.NewV7()

	_, err := s.db.Exec(`
	insert into categories (id, name)
		values (?, ?)
	`, category.ID.String(), category.Name)

	if err != nil {
		var sqlErr *sqlite.Error
		if errors.As(err, &sqlErr) {
			if sqlErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return nil, model.ErrCategoryAlreadyExists
			}
		}
		return nil, err
	}

	return &category, nil
}

func (s *SQLiteStore) ListCategories() ([]model.Category, error) {
	rows, err := s.db.Query(`
		select id, name
		from categories
		order by name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)

	var category model.Category

	for rows.Next() {
		err = rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *SQLiteStore) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	row := s.db.QueryRow(`
		select id, name
		from categories
		where id = ?
	`, id)

	var category model.Category
	err := row.Scan(
		&category.ID,
		&category.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil
}
