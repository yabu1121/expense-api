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
	dsn := filePath + "?_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
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

	if err := sqliteStore.createTagTable(); err != nil {
		db.Close()
		return nil, err
	}

	if err := sqliteStore.createExpenseTagTable(); err != nil {
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
			id TEXT NOT NULL,
			name TEXT NOT NULL UNIQUE,

			PRIMARY KEY (id)
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
			id TEXT NOT NULL,
			title TEXT NOT NULL,
			amount INTEGER NOT NULL,
			category_id TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		)`,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) createTagTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id TEXT NOT NULL,
			name TEXT NOT NULL UNIQUE,

			PRIMARY KEY (id)
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) createExpenseTagTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS expense_tags (
			expense_id TEXT NOT NULL,
			tag_id TEXT NOT NULL,

			PRIMARY KEY (expense_id, tag_id),
			FOREIGN KEY (expense_id) REFERENCES expenses(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) AddTagToExpense(expenseID, tagID uuid.UUID) error {
	_, err := s.db.Exec(`
		insert into expense_tags (expense_id, tag_id) values (?, ?)
	`, expenseID.String(), tagID.String())
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) ListExpenses() ([]model.Expense, error) {
	rows, err := s.db.Query(`
		select id, title, amount, category_id, created_at
		from expenses
		order by created_at desc
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
			&expense.CategoryID,
			&expense.CreatedAt,
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

func (s *SQLiteStore) GetExpenseByID(id uuid.UUID) (*model.Expense, error) {
	result := s.db.QueryRow(`
		select id, title, amount, category_id, created_at
		from expenses
		where id = ?
	`, id.String())

	var expense model.Expense
	err := result.Scan(
		&expense.ID,
		&expense.Title,
		&expense.Amount,
		&expense.CategoryID,
		&expense.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrExpenseNotFound
		}
		return nil, err
	}

	return &expense, nil
}

func (s *SQLiteStore) CreateExpense(expense model.ExpenseRequest) (*model.Expense, error) {
	id := uuid.NewV7()
	result := s.db.QueryRow(`
		insert into expenses (id, title, amount, category_id)
		values(?, ?, ?, ?)
		returning id, title, amount, category_id, created_at
	`, id.String(), expense.Title, expense.Amount, expense.CategoryID.String())

	var createdExpense model.Expense
	err := result.Scan(
		&createdExpense.ID,
		&createdExpense.Title,
		&createdExpense.Amount,
		&createdExpense.CategoryID,
		&createdExpense.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdExpense, nil
}

func (s *SQLiteStore) UpdateExpense(id uuid.UUID, expense model.ExpenseRequest) (*model.Expense, error) {
	result := s.db.QueryRow(`
		update expenses
		set title = ?, amount = ?, category_id = ?
		where id = ?
		returning id, title, amount, category_id, created_at
	`, expense.Title, expense.Amount, expense.CategoryID.String(), id.String())

	var updatedExpense model.Expense
	err := result.Scan(
		&updatedExpense.ID,
		&updatedExpense.Title,
		&updatedExpense.Amount,
		&updatedExpense.CategoryID,
		&updatedExpense.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrExpenseNotFound
		}
		return nil, err
	}

	return &updatedExpense, nil
}

func (s *SQLiteStore) DeleteExpense(id uuid.UUID) error {
	res, err := s.db.Exec(`
		delete from expenses
		where id = ?
	`, id.String())
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
	if err := row.Scan(&summary.Count, &summary.TotalAmount); err != nil {
		return nil, err
	}
	return &summary, nil
}

func (s *SQLiteStore) CreateCategory(req model.CategoryRequest) (*model.Category, error) {
	category := model.Category{
		Name: req.Name,
	}

	category.ID = uuid.NewV7()

	_, err := s.db.Exec(`
	insert into categories (id, name) values (?, ?)
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
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
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
	`, id.String())

	var category model.Category
	if err := row.Scan(&category.ID, &category.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (s *SQLiteStore) ListTags() ([]model.Tag, error) {
	rows, err := s.db.Query(`
		select id, name
		from tags
		order by name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]model.Tag, 0)
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *SQLiteStore) CreateTag(tag model.Tag) (*model.Tag, error) {
	tag.ID = uuid.NewV7()

	_, err := s.db.Exec(`
		insert into tags (id, name) values (?, ?)
	`, tag.ID.String(), tag.Name)
	if err != nil {
		var sqlErr *sqlite.Error
		if errors.As(err, &sqlErr) && sqlErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil, model.ErrTagAlreadyExists
		}
		return nil, err
	}
	return &tag, nil
}
