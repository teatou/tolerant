package relational

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(host string, port int, user, password, dbName string) (*Storage, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, fmt.Errorf("database connetction: %w", err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Add(to, sum int) error {
	q := `INSERT INTO users (id, sum) 
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE 
	  	SET sum = excluded.sum + $2;`

	_, err := s.db.Exec(q, to, sum) // ADD TX!!
	return err
}

func (s *Storage) Retrieve(from, sum int) error {
	curSum, err := s.CheckBalance(from)
	if err != nil {
		return err
	}

	if curSum >= sum {
		q := `UPDATE users
		SET sum = sum - $2
		WHERE id = $1`

		_, err = s.db.Exec(q, sum)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not enough money")
	}

	return nil
}

func (s *Storage) Transfer(from, to, sum int) error {
	// !имплементировать откаты!

	sumFrom, err := s.CheckBalance(from)
	if err != nil {
		return err
	}

	if sumFrom < sum {
		return fmt.Errorf("not enough money")
	}

	err = s.Add(to, sum)
	if err != nil {
		return err
	}

	err = s.Retrieve(from, sum)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CheckBalance(id int) (int, error) {
	q := `SELECT sum FROM users
		WHERE id = $1`

	var sum int
	if err := s.db.QueryRow(q, id).Scan(&sum); err != nil {
		return 0, err
	}

	return sum, nil
}
