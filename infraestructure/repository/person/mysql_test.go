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
	"testing"
)

func TestCRUD(t *testing.T) {
	ctx := context.Background()
	container, err := person.SetupMysqL(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)
	db, err := sql.Open("mysql", container.URI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	err = person.InitMySQL(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	defer person.TruncateMySQL(ctx, db)

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
	container, err := person.SetupMysqL(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)
	db, err := sql.Open("mysql", container.URI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	err = person.InitMySQL(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	defer person.TruncateMySQL(ctx, db)

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
