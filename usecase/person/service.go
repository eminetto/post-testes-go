package person

import (
	"fmt"
	"github.com/PicPay/go-test-workshop/entity"
)

type Service struct {
	r Repository
}

//NewService cria um novo serviço. Lembre-se: receba interfaces, retorne structs ;)
func NewService(r Repository) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) Get(id entity.ID) (*entity.Person, error) {
	p, err := s.r.Get(id)
	if err != nil {
		return nil, fmt.Errorf("erro lendo person do repositório: %w", err)
	}
	return p, nil
}

func (s *Service) Search(query string) ([]*entity.Person, error) {
	p, err := s.r.Search(query)
	if err != nil {
		return nil, fmt.Errorf("erro buscando person do repositório: %w", err)
	}
	return p, nil
}

func (s *Service) List() ([]*entity.Person, error) {
	p, err := s.r.List()
	if err != nil {
		return nil, fmt.Errorf("erro listando person do repositório: %w", err)
	}
	return p, nil
}

func (s *Service) Create(firstName, lastName string) (entity.ID, error) {
	p := entity.Person{
		Name:     firstName,
		LastName: lastName,
	}
	id, err := s.r.Create(&p)
	if err != nil {
		return 0, fmt.Errorf("erro criando person no repositório: %w", err)
	}
	return id, nil
}

func (s *Service) Update(e *entity.Person) error {
	err := s.r.Update(e)
	if err != nil {
		return fmt.Errorf("erro atualizando person no repositório: %w", err)
	}
	return nil
}

func (s *Service) Delete(id entity.ID) error {
	err := s.r.Delete(id)
	if err != nil {
		return fmt.Errorf("erro removendo person do repositório: %w", err)
	}
	return nil
}
