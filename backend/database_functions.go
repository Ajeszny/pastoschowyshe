package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
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
	query = fmt.Sprintf("CREATE DATABASE $1")
	stmt, err = conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(os.Getenv("DB_NAME"))
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

	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS pasty (pasta_id SERIAL PRIMARY KEY, pasta_name varchar(100), pasta_body varchar(100000))")
	stmt, err := connection.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	stmt.Close()
	if err != nil {
		return err
	}

	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS pasty_tag_relation (pasta_id integer, tag_id integer)")
	stmt, err = connection.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	stmt.Close()
	if err != nil {
		return err
	}

	query = fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS Tags (tag_id SERIAL PRIMARY KEY, tag_root varchar(100), tag_name varchar(130))")
	stmt, err = connection.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func add_new_record(new pasta) error {
	query := fmt.Sprintf("INSERT INTO pasty (pasta_name, pasta_body)  VALUES($1, $2)")
	stmt, err := connection.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(new.Name, new.Text)
	if err != nil {
		return err
	}
	stmt.Close()
	query = "SELECT pasta_id FROM pasty WHERE pasta_name = $1"
	stmt, err = connection.Prepare(query)
	if err != nil {
		return err
	}
	qres, err := stmt.Query(new.Name)
	var pasta_id int
	if qres.Next() {
		qres.Scan(&pasta_id)
	}
	if err != nil {
		return err
	}
	stmt.Close()

	for _, tag := range new.Tags {
		query = "SELECT tag_id FROM Tags WHERE tag_root = $1"
		stmt, err = connection.Prepare(query)
		if err != nil {
			return err
		}
		result, err := stmt.Query(tag)
		var tag_id int
		if result.Next() {
			result.Scan(&tag_id)
		}
		if err != nil {
			return err
		}
		stmt.Close()
		query = "INSERT INTO pasty_tag_relation VALUES($1, $2)"
		stmt, err := connection.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(pasta_id, tag_id)
		if err != nil {
			return err
		}
		stmt.Close()
	}

	return nil
}

func get_records() ([]pasta, error) {
	query := "SELECT pasta_id, pasta_name FROM pasty"
	stmt, err := connection.Prepare(query)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Query()
	retval := []pasta{}
	i := 0
	for result.Next() {
		retval = append(retval, pasta{})
		result.Scan(&retval[i].Id, &retval[i].Name)
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
