package person

import "github.com/PicPay/go-test-workshop/entity"

type Reader interface {
	Get(id entity.ID) (*entity.Person, error)
	Search(query string) ([]*entity.Person, error)
	List() ([]*entity.Person, error)
}

type Writer interface {
	Create(e *entity.Person) (entity.ID, error)
	Update(e *entity.Person) error
	Delete(id entity.ID) error
}

type Repository interface {
	Reader
	Writer
}

/*
É uma boa prática da comunidade Go criarmos interfaces pequenas e usarmos composição
Desta forma poderíamos ter um repositório que só implementa a leitura ou a gravação
*/

/*
Outra boa prática. O definição da interface é feita do lado de quem a usa, e não de quem a implementa
Desta forma, como o UseCase precisa do Repository, a definição da interface fica aqui enquanto que a
implementação vai ser feita no pacote infrastructure/repository/person
*/

type UseCase interface {
	Get(id entity.ID) (*entity.Person, error)
	Search(query string) ([]*entity.Person, error)
	List() ([]*entity.Person, error)
	Create(firstName, lastName string) (entity.ID, error)
	Update(e *entity.Person) error
	Delete(id entity.ID) error
}
