package person

import (
	"database/sql"
	"fmt"
	"github.com/PicPay/go-test-workshop/entity"
	"strings"
	"time"
)

//MySQL mysql repo
type MySQL struct {
	db *sql.DB
}

//NewMySQL create new repository
func NewMySQL(db *sql.DB) *MySQL {
	return &MySQL{
		db: db,
	}
}

//Create a person
func (r *MySQL) Create(p *entity.Person) (entity.ID, error) {
	stmt, err := r.db.Prepare(`
		insert into person (first_name, last_name, created_at) 
		values(?,?,?)`)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(
		p.Name,
		p.LastName,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return entity.ID(id), nil
}

//Get a person
func (r *MySQL) Get(id entity.ID) (*entity.Person, error) {
	stmt, err := r.db.Prepare(`select id, first_name, last_name from person where id = ?`)
	if err != nil {
		return nil, err
	}
	var p entity.Person
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, fmt.Errorf("not found")
	}
	err = rows.Scan(&p.ID, &p.Name, &p.LastName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

//Update a person
func (r *MySQL) Update(p *entity.Person) error {
	_, err := r.db.Exec("update person set first_name = ?, last_name = ?, updated_at = ? where id = ?", p.Name, p.LastName, time.Now().Format("2006-01-02"), p.ID)
	if err != nil {
		return err
	}
	return nil
}

//Search person
func (r *MySQL) Search(query string) ([]*entity.Person, error) {
	stmt, err := r.db.Prepare(`select id, first_name, last_name from person where first_name like ? or last_name like ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var people []*entity.Person
	query = "%" + strings.ToLower(query) + "%"
	rows, err := stmt.Query(query, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p entity.Person
		err = rows.Scan(&p.ID, &p.Name, &p.LastName)
		if err != nil {
			return nil, err
		}
		people = append(people, &p)
	}
	if len(people) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return people, nil
}

//List person
func (r *MySQL) List() ([]*entity.Person, error) {
	stmt, err := r.db.Prepare(`select id, first_name, last_name from person`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var people []*entity.Person
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p entity.Person
		err = rows.Scan(&p.ID, &p.Name, &p.LastName)
		if err != nil {
			return nil, err
		}
		people = append(people, &p)
	}
	if len(people) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return people, nil
}

//Delete a person
func (r *MySQL) Delete(id entity.ID) error {
	p, _ := r.Get(id)
	if p == nil {
		return fmt.Errorf("not found")
	}
	_, err := r.db.Exec("delete from person where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
