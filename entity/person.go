package entity

//ID representa o ID de uma entidade.
//É uma boa prática criarmos esse tipo, pois se em algum momento precisarmos mudar para outro formato (UUID por exemplo)
//não quebramos o restante do projeto
type ID int

type Person struct {
	ID       ID
	Name     string
	LastName string
}
