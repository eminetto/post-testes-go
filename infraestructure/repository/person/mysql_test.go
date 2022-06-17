//go:build integration

package person_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/PicPay/go-test-workshop/entity"
	"github.com/PicPay/go-test-workshop/infraestructure/repository/person"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

type mysqlDBContainer struct {
	testcontainers.Container
	URI string
}

const (
	dbUser         = "workshop"
	dbPassword     = "workshop"
	database       = "workshop"
	dbRootPassword = "db-root-password"
)

func setupMysqL(ctx context.Context) (*mysqlDBContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mariadb:latest",
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor:   wait.ForLog("Version: '10.8.3-MariaDB-1:10.8.3+maria~jammy'  socket: '/run/mysqld/mysqld.sock'  port: 3306  mariadb.org binary distribution"),
		Env: map[string]string{
			"MARIADB_USER":          dbUser,
			"MARIADB_PASSWORD":      dbPassword,
			"MARIADB_ROOT_PASSWORD": dbRootPassword,
			"MARIADB_DATABASE":      database,
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	mappedPort, err := container.MappedPort(ctx, "3306")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", "root", dbRootPassword, hostIP, mappedPort.Port(), database)

	return &mysqlDBContainer{Container: container, URI: uri}, nil
}

func initMySQL(ctx context.Context, db *sql.DB) error {
	query := []string{
		fmt.Sprintf("use %s;", database),
		"create table if not exists person (id int AUTO_INCREMENT,first_name varchar(100), last_name varchar(100), created_at datetime, updated_at datetime, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;",
	}
	for _, q := range query {
		_, err := db.ExecContext(ctx, q)
		if err != nil {
			return err
		}
	}

	return nil
}

func truncateMySQL(ctx context.Context, db *sql.DB) error {
	query := []string{
		fmt.Sprintf("use %s;", database),
		"truncate table person",
	}
	for _, q := range query {
		_, err := db.ExecContext(ctx, q)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestCRUD(t *testing.T) {
	ctx := context.Background()
	container, err := setupMysqL(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)
	db, err := sql.Open("mysql", container.URI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	err = initMySQL(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	defer truncateMySQL(ctx, db)

	repo := person.NewMySQL(db)

	t.Run("inserir person", func(t *testing.T) {
		p := &entity.Person{
			Name:     "Ozzy",
			LastName: "Osbourne",
		}
		id, err := repo.Create(p)
		assert.Equal(t, entity.ID(1), id)
		assert.Nil(t, err)
	})
	t.Run("recuperar person", func(t *testing.T) {
		result, err := repo.Get(entity.ID(1))
		assert.Equal(t, "Ozzy", result.Name)
		assert.Nil(t, err)
	})
	t.Run("atualizar person", func(t *testing.T) {
		result, err := repo.Get(entity.ID(1))
		assert.Nil(t, err)
		result.Name = "Novo nome"
		err = repo.Update(result)
		assert.Nil(t, err)
		saved, err := repo.Get(entity.ID(1))
		assert.Nil(t, err)
		assert.Equal(t, "Novo nome", saved.Name)
	})
	t.Run("listar person", func(t *testing.T) {
		result, err := repo.List()
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Osbourne", result[0].LastName)
		assert.Nil(t, err)
	})
	t.Run("remover person", func(t *testing.T) {
		err := repo.Delete(entity.ID(1))
		assert.Nil(t, err)
	})
	t.Run("listar person vazia", func(t *testing.T) {
		result, err := repo.List()
		assert.Nil(t, result)
		assert.Errorf(t, err, "not found")
	})
	t.Run("remover person não existente", func(t *testing.T) {
		err := repo.Delete(entity.ID(1))
		assert.Errorf(t, err, "not found")
	})
}

func TestSearch(t *testing.T) {
	ctx := context.Background()
	container, err := setupMysqL(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)
	db, err := sql.Open("mysql", container.URI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	err = initMySQL(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	defer truncateMySQL(ctx, db)

	repo := person.NewMySQL(db)

	p1 := &entity.Person{
		Name:     "Ozzy",
		LastName: "Osbourne",
	}
	p2 := &entity.Person{
		Name:     "Ronnie",
		LastName: "Dio",
	}
	p1.ID, err = repo.Create(p1)
	assert.Nil(t, err)
	p2.ID, err = repo.Create(p2)
	assert.Nil(t, err)

	tests := []struct {
		query       string
		result      []*entity.Person
		expectedErr error
	}{
		{
			query:       "ozzy",
			result:      []*entity.Person{p1},
			expectedErr: nil,
		},
		{
			query:       "Ozzy",
			result:      []*entity.Person{p1},
			expectedErr: nil,
		},
		{
			query:       "osbourne",
			result:      []*entity.Person{p1},
			expectedErr: nil,
		},
		{
			query:       "Osbourne",
			result:      []*entity.Person{p1},
			expectedErr: nil,
		},
		{
			query:       "Dio",
			result:      []*entity.Person{p2},
			expectedErr: nil,
		},
		{
			query:       "dio",
			result:      []*entity.Person{p2},
			expectedErr: nil,
		},
		{
			query:       "ronnie",
			result:      []*entity.Person{p2},
			expectedErr: nil,
		},
		{
			query:       "Ronnie",
			result:      []*entity.Person{p2},
			expectedErr: nil,
		},
		{
			query:       "Tony",
			result:      nil,
			expectedErr: fmt.Errorf("not found"),
		},
		{
			query:       "martin",
			result:      nil,
			expectedErr: fmt.Errorf("not found"),
		},
	}
	for _, test := range tests {
		found, err := repo.Search(test.query)
		assert.Equal(t, test.expectedErr, err)
		assert.Equal(t, test.result, found)
	}

}
