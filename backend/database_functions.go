package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"os"
	"strings"
)

var connection *sql.DB

func create_db() error {
	conn, err := sql.Open("postgres",
		fmt.Sprintf("host=localhost user=%s password=%s port=5432 dbname=postgres sslmode=disable",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")))
	if err != nil {
		return err
	}
	query := fmt.Sprintf("SELECT FROM pg_database WHERE datname = 'pastoshowyshe'")
	stmt, err := conn.Prepare(query)
	if err != nil {
		return err
	}
	rows, err := stmt.Query()
	if rows.Next() {
		return nil
	}
	stmt.Close()
	query = fmt.Sprintf("CREATE DATABASE pastoshowyshe")
	stmt, err = conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func populate_db() error {
	conn, err := sql.Open("postgres",
		fmt.Sprintf("host=localhost user=%s password=%s port=5432 dbname=%s sslmode=disable",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
	if err != nil {
		return err
	}
	connection = conn
	query, err := os.ReadFile(os.Getenv("DB_FILE_PATH"))
	if err != nil {
		return err
	}
	_, err = connection.Exec(string(query))
	if err != nil {
		return err
	}
	return nil
}

func add_new_record(new pasta) error {
	query := fmt.Sprintf("SELECT insert_pasta($1, $2, $3)")
	stmt, err := connection.Prepare(query)
	if err != nil {
		return err
	}
	tags := "{" + strings.Join(new.Tags, ",") + "}"
	_, err = stmt.Query(new.Name, new.Text, tags)
	if err != nil {
		return err
	}
	return nil
}

func get_records() ([]pasta, error) {
	query := "SELECT pastaid, pastaname, tags from get_records()"
	stmt, err := connection.Prepare(query)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Query()
	retval := []pasta{}

	i := 0
	for result.Next() {
		retval = append(retval, pasta{})
		var tags []sql.NullString
		err = result.Scan(&retval[i].Id, &retval[i].Name, pq.Array(&tags))
		if err != nil {
			return retval, err
		}
		for _, tag := range tags {
			if tag.Valid {
				retval[i].Tags = append(retval[i].Tags, tag.String)
			}
		}
		i++
	}
	//TODO: add tag support
	return retval, nil
}

func get_pasta(id int) (pasta, error) {
	query := "SELECT pasta_id, pasta_name, pasta_body FROM pasty WHERE pasta_id = $1"
	stmt, err := connection.Prepare(query)
	var p pasta
	if err != nil {
		return p, err
	}
	result, err := stmt.Query(id)
	if !result.Next() {
		return p, errors.New("pasta does not exist")
	}
	err = result.Scan(&p.Id, &p.Name, &p.Text)
	if err != nil {
		return p, err
	}
	return p, nil
}
