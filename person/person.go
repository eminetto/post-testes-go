package person

//ID representa o ID de uma entidade.
//É uma boa prática criarmos esse tipo, pois se em algum momento precisarmos mudar para outro formato (UUID por exemplo)
//não quebramos o restante do projeto
type ID int

//Person define o que é uma pessoa
type Person struct {
	ID       ID
	Name     string
	LastName string
}

type Reader interface {
	Get(id ID) (*Person, error)
	Search(query string) ([]*Person, error)
	List() ([]*Person, error)
}

type Writer interface {
	Create(e *Person) (ID, error)
	Update(e *Person) error
	Delete(id ID) error
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
implementação vai ser feita no pacote person/mysql
*/

type UseCase interface {
	Get(id ID) (*Person, error)
	Search(query string) ([]*Person, error)
	List() ([]*Person, error)
	Create(firstName, lastName string) (ID, error)
	Update(e *Person) error
	Delete(id ID) error
}
