//go:build integration

package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/PicPay/go-test-workshop/person"
	"github.com/PicPay/go-test-workshop/person/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PersonTestSuite struct {
	suite.Suite
	ctx       context.Context
	container *person.MysqlDBContainer
	db        *sql.DB
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPersonTestSuite(t *testing.T) {
	suite.Run(t, new(PersonTestSuite))
}

// before each test
func (suite *PersonTestSuite) SetupTest() {
	var err error
	suite.ctx = context.Background()
	suite.container, err = person.SetupMysqL(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.db, err = sql.Open("mysql", suite.container.URI)
	if err != nil {
		suite.T().Error(err)
	}
	err = person.InitMySQL(suite.ctx, suite.db)
	if err != nil {
		suite.T().Fatal(err)
	}
}

//This method is ran after all tests have runned, it cleans the suite
func (suite *PersonTestSuite) TearDownTest() {
	person.TruncateMySQL(suite.ctx, suite.db)
	err := suite.db.Close()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.container.Terminate(suite.ctx)
}

func (suite *PersonTestSuite) TestCRUD() {
	repo := mysql.NewMySQL(suite.db)

	suite.T().Run("inserir person", func(t *testing.T) {
		p := &person.Person{
			Name:     "Ozzy",
			LastName: "Osbourne",
		}
		id, err := repo.Create(p)
		assert.Equal(t, person.ID(1), id)
		assert.Nil(t, err)
	})
	suite.T().Run("recuperar person", func(t *testing.T) {
		result, err := repo.Get(person.ID(1))
		assert.Equal(t, "Ozzy", result.Name)
		assert.Nil(t, err)
	})
	suite.T().Run("atualizar person", func(t *testing.T) {
		result, err := repo.Get(person.ID(1))
		assert.Nil(t, err)
		result.Name = "Novo nome"
		err = repo.Update(result)
		assert.Nil(t, err)
		saved, err := repo.Get(person.ID(1))
		assert.Nil(t, err)
		assert.Equal(t, "Novo nome", saved.Name)
	})
	suite.T().Run("listar person", func(t *testing.T) {
		result, err := repo.List()
		assert.Equal(suite.T(), 1, len(result))
		assert.Equal(suite.T(), "Osbourne", result[0].LastName)
		assert.Nil(suite.T(), err)
	})
	suite.T().Run("remover person", func(t *testing.T) {
		err := repo.Delete(person.ID(1))
		assert.Nil(t, err)
	})
	suite.T().Run("listar person vazia", func(t *testing.T) {
		result, err := repo.List()
		assert.Nil(suite.T(), result)
		assert.Errorf(suite.T(), err, "not found")
	})
	suite.T().Run("remover person não existente", func(t *testing.T) {
		err := repo.Delete(person.ID(1))
		assert.Errorf(suite.T(), err, "not found")
	})
}
func (suite *PersonTestSuite) TestSearch() {
	repo := mysql.NewMySQL(suite.db)

	p1 := &person.Person{
		Name:     "Ozzy",
		LastName: "Osbourne",
	}
	p2 := &person.Person{
		Name:     "Ronnie",
		LastName: "Dio",
	}
	var err error
	p1.ID, err = repo.Create(p1)
	assert.Nil(suite.T(), err)
	p2.ID, err = repo.Create(p2)
	assert.Nil(suite.T(), err)

	tests := []struct {
		query       string
		result      []*person.Person
		expectedErr error
	}{
		{
			query:       "ozzy",
			result:      []*person.Person{p1},
			expectedErr: nil,
		},
		{
			query:       "Ozzy",
			result:      []*person.Person{p1},
			expectedErr: nil,
		},
		{
			query:       "osbourne",
			result:      []*person.Person{p1},
			expectedErr: nil,
		},
		{
			query:       "Osbourne",
			result:      []*person.Person{p1},
			expectedErr: nil,
		},
		{
			query:       "Dio",
			result:      []*person.Person{p2},
			expectedErr: nil,
		},
		{
			query:       "dio",
			result:      []*person.Person{p2},
			expectedErr: nil,
		},
		{
			query:       "ronnie",
			result:      []*person.Person{p2},
			expectedErr: nil,
		},
		{
			query:       "Ronnie",
			result:      []*person.Person{p2},
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
		assert.Equal(suite.T(), test.expectedErr, err)
		assert.Equal(suite.T(), test.result, found)
	}

}
