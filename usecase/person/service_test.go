//go:build unit

package person_test

//boa prática: criar um pacote _test para que sejam testadas as funções públicas do pacote e não as internas

import (
	"fmt"
	"github.com/PicPay/go-test-workshop/entity"
	"github.com/PicPay/go-test-workshop/usecase/person"
	"github.com/PicPay/go-test-workshop/usecase/person/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Get(t *testing.T) {
	t.Run("usuário encontrado", func(t *testing.T) {
		//fase: Configure os dados de teste
		p := &entity.Person{
			ID:       1,
			Name:     "Ozzy",
			LastName: "Osbourne",
		}
		repo := mocks.NewRepository(t)
		repo.On("Get", entity.ID(1)).
			Return(p, nil).
			Once()
		service := person.NewService(repo)
		//fase: Invoque o método sendo testado
		found, err := service.Get(entity.ID(1))

		//fase: Confirme que os resultados esperados são retornados
		assert.Nil(t, err)
		assert.Equal(t, p, found)

	})
	t.Run("usuário não encontrado", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.On("Get", entity.ID(1)).
			Return(nil, fmt.Errorf("not found")).
			Once()
		service := person.NewService(repo)
		found, err := service.Get(entity.ID(1))
		assert.Nil(t, found)
		assert.Errorf(t, err, "erro lendo person do repositório: %w")
	})
}

func TestService_Search(t *testing.T) {
	//aqui vamos usar uma técnica chamada Table based tests
	p1 := &entity.Person{
		ID:       1,
		Name:     "Ozzy",
		LastName: "Osbourne",
	}
	p2 := &entity.Person{
		ID:       2,
		Name:     "Ronnie",
		LastName: "Dio",
	}

	tests := []struct {
		query       string
		result      []*entity.Person
		expectedErr error
		mockErr     error
	}{
		{
			query:       "ozzy",
			result:      []*entity.Person{p1},
			expectedErr: nil,
			mockErr:     nil,
		},
		{
			query:       "Ozzy",
			result:      []*entity.Person{p1},
			expectedErr: nil,
			mockErr:     nil,
		},
		{
			query:       "osbourne",
			result:      []*entity.Person{p1},
			expectedErr: nil,
			mockErr:     nil,
		},
		{
			query:       "Osbourne",
			result:      []*entity.Person{p1},
			expectedErr: nil,
			mockErr:     nil,
		},
		{
			query:       "Dio",
			result:      []*entity.Person{p2},
			expectedErr: nil,
			mockErr:     nil,
		},
		{
			query:       "dio",
			result:      []*entity.Person{p2},
			expectedErr: nil,
			mockErr:     nil,
		},
		{
			query:       "ronnie",
			result:      []*entity.Person{p2},
			expectedErr: nil,
			mockErr:     nil,
		},
		{
			query:       "Ronnie",
			result:      []*entity.Person{p2},
			expectedErr: nil,
			mockErr:     nil,
		},
		{
			query:       "Tony",
			result:      nil,
			expectedErr: fmt.Errorf("erro buscando person do repositório: %w", fmt.Errorf("not found")),
			mockErr:     fmt.Errorf("not found"),
		},
		{
			query:       "martin",
			result:      nil,
			expectedErr: fmt.Errorf("erro buscando person do repositório: %w", fmt.Errorf("not found")),
			mockErr:     fmt.Errorf("not found"),
		},
	}
	for _, test := range tests {
		repo := mocks.NewRepository(t)
		repo.On("Search", test.query).
			Return(test.result, test.mockErr).
			Once()
		service := person.NewService(repo)
		found, err := service.Search(test.query)

		assert.Equal(t, test.expectedErr, err)
		assert.Equal(t, test.result, found)
	}

}

//para fins didáticos, deixo os demais testes para serem implementados como aprendizado ;)
