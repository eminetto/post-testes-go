# Testes automatizados em Go

Repositório criado para exemplificar os tipos de testes e boas práticas que podem ser aplicados. Este código foi usado como exemplo neste [post](https://medium.com/inside-picpay/testes-automatizados-em-go-aa5cf9ed672e)

Baseado na [aplicação](https://github.com/hamvocke/spring-testing) usada como exemplo.


## Arquitetura da aplicação

```
 ╭┄┄┄┄┄┄┄╮      ┌──────────┐      ┌──────────┐
 ┆   ☁   ┆  ←→  │          │  ←→  │    💾    │
 ┆  Web  ┆ HTTP │    Go    │      │ Database │
 ╰┄┄┄┄┄┄┄╯      │  Service │      └──────────┘
                └──────────┘
                     ↑ JSON/HTTP
                     ↓
                ┌──────────┐
                │    ☁     │
                │ Weather  │
                │   API    │
                └──────────┘

```

A aplicação fornece três endpoints:

```

GET /hello: Retorna "Hello World!". 
GET /hello/{lastname}: Procura no banco de dados a pessoa pelo seu sobrenome e retorna "Hello {Firstname} {Lastname}" se a pessoa é encontrada. Retorna 404 caso não encontrada.
GET /weather/{lat}/{long}: Chama uma API de previsão do tempo via HTTP e retorna as condições de acordo com as coordenadas. Retorna 404 caso não encontrada.

```

## Testes


Neste repositório podemos ver implementações da Pirâmide de Testes

```
          /\
         /  \
        /    \  End to end
       /      \ 
      /────────\
     /          \  Integração
    /            \
   /──────────────\
  /                \
 /                  \ Unitários 
/                    \
 ────────────────────

```

### Estrutura dos testes

Antes de mergulhar nos tipos de teste,  uma boa estrutura para todos os testes é esta:

1. Configure os dados de teste, prepare o teste
2. Invoque o método/função sendo testada, execute o teste
3. Confirme que os resultados esperados são retornados, verifique as asserções

Este padrão também é conhecido como *Arrange* (Prepare o teste), *Act* (Execute o teste) e *Assert* (Verifique as asserções). Vamos observar esta estrutura em todos os testes.

### Testes unitários

Testes de unidade garantem que uma determinada unidade (o *sujeito em teste*) da base de código funcione conforme o esperado. Os testes de unidade têm o escopo mais restrito de todos os testes do conjunto de testes. O número de testes de unidade do conjunto de testes superará em grande parte qualquer outro tipo de teste.


#### O que testar?

Os testes unitários devem pelo menos testar a interface pública do pacote.  Em Go é possível testar tanto as funções públicas (as que começam com a primeira letra maiúscula) quanto as funções privadas do pacote, mas é recomendado testarmos prioritariamente as públicas.  

Há uma linha tênue quando se trata de escrever testes de unidade: eles devem garantir que todos os seus caminhos de código não triviais sejam testados (incluindo caminho feliz e casos de borda). Ao mesmo tempo, eles não devem estar muito vinculados à sua implementação.

Por que isso?

Testes muito próximos do código de produção rapidamente se tornam irritantes. Assim que você refatorar seu código de produção (recapitulação rápida: refatorar significa alterar a estrutura interna do seu código sem alterar o comportamento visível externamente), seus testes de unidade irão quebrar.
Resumindo, não reflita sua estrutura de código interna em seus testes de unidade. Teste para comportamento observável em vez disso. Para ilustrar esse conceito, no código a seguir:


Arquivo `service.go`
```go

package service

//NewService create new service
func NewService() *Service {
    return &Service{}
}

//FindAll
func (s *Service) FindAll() ([]*entity.Privilege, error) {
    acl := s.getDefaultPrivileges()
    return acl, nil
}

//FindByRole
func (s *Service) FindyByRole(r *entity.Role) ([]*entity.Privilege, error) {
    var ret []*entity.Privilege
    acl := s.getDefaultPrivileges()
    for _, p := range acl {
        if p.Role.Slug == r.Slug {
            ret = append(ret, p)
        }
        pChildren := walkPrivilegeChildren(p)
        for _, pC := range pChildren {
            if pC.Role.Slug == r.Slug {
                ret = append(ret, pC)
            }
        }
    }
    return ret, nil
}

func walkPrivilegeChildren(priv *entity.Privilege) []*entity.Privilege {
    var p []*entity.Privilege
    for _, c := range priv.Children {
        p = append(p, walkPrivilegeChildren(c)...)
    }
    return p
}



```

O recomendado é criarmos testes para a interface pública do pacote, as funções `NewService`, `FindAll` e `FindyByRole`. Desta forma,
se for necessário uma refatoração nas funções internas, como a `getDefaultPrivileges` e `walkPrivilegeChildren` não é necessário refatorar também os testes unitários. 
Para fazer isso em Go basta criar um pacote especial no momento da escrita do teste:

Arquivo `service_test.go`

```go
package service_test

import (
	"testing"
	"github.com/PicPay/example"
)

func TestFindAll(t *testing.T) {
	s := example.NewService()
	all, err := s.FindAll()
	//asserts vão aqui
}
```

Desta forma, nosso teste se comporta como um pacote diferente, apesar do arquivo estar no mesmo diretório que o `service.go`. Essa é uma facilidade da linguagem para facilitar a criação de testes.


#### Exemplos de teste unitário

[person/service_test.go](https://github.com/eminetto/post-testes-go/blob/main/person/service_test.go)

Este arquivo contém os testes do serviço que implementa a interface [UseCase](https://github.com/eminetto/post-testes-go/blob/main/person/person.go#L43). 

Como o serviço tem por dependência uma implementação da interface [Repository](https://github.com/eminetto/post-testes-go/blob/main/person/person.go#L27) (que por sua vez precisa de uma conexão com o banco de dados), vamos usar o conceito de [mocks](https://martinfowler.com/articles/mocksArentStubs.html) para mantermos o foco do teste apenas na regra de negócio do serviço.
Para gerarmos facilmente os `mocks` estamos usando a ferramenta [mockery](https://github.com/vektra/mockery), que lê as interfaces e gera código para usarmos nos testes.
A geração dos `mocks` é executada pelo comando `make generate-mocks` e pode ser executada manualmente ou automaticamente quando executamos o comando `make unit-test`


[weather/service_test.go](https://github.com/eminetto/post-testes-go/blob/main/weather/service_test.go)

Este arquivo contém os testes do serviço que implementa a interface [UseCase](https://github.com/eminetto/post-testes-go/blob/main/weather/weather.go#L37).

Este é um serviço que faz uso de uma [API externa](https://api.openweathermap.org/). Para não acessar a API real a cada teste criamos um `mock` para simular o seu comportamento. 
Vale destacar uma boa prática neste pacote. Ao invés de colocarmos como dependência do `UseCase` um `http.Client` padrão da linguagem foi criada uma [interface](https://github.com/eminetto/post-testes-go/blob/main/weather/weather.go#L33) para ser usada como dependência. 
No [construtor](https://github.com/eminetto/post-testes-go/blob/main/weather/service.go#L19) do serviço criamos uma instância de `http.Client` e damos a opção do usuário substituir esse cliente padrão por outra implementação. 
Fazemos uso desta opção no [momento do teste](https://github.com/eminetto/post-testes-go/blob/main/weather/service_test.go#L30) ao passar um `mock` do client. 
Esta implementação pode ser resumida pela frase `“Don’t Mock What You Don’t Own”` e mais detalhes podem ser vistos neste [post](https://hynek.me/articles/what-to-mock-in-5-mins/).


[internal/http/echo/handler_test.go](https://github.com/eminetto/post-testes-go/blob/main/internal/http/echo/handler_test.go)

Neste arquivo implementamos os testes unitários da camada de API. 

Eles usam os `mocks` da camada de `UseCase`. 

#### Executando os testes unitários

Execute

    make unit-test


### Testes de integração

Todos os aplicativos não triviais serão integrados com algumas outras partes (bancos de dados, sistemas de arquivos, chamadas de rede para outros aplicativos). 
Ao escrever testes de unidade, essas são geralmente as partes que você deixa de fora para obter um melhor isolamento e testes mais rápidos. 
Ainda assim, seu aplicativo irá interagir com outras partes e isso precisa ser testado. 

[Testes de integração](https://martinfowler.com/bliki/IntegrationTest.html) estão disponíveis para ajudar. 
**Eles testam a integração do seu aplicativo com todas as partes que vivem fora do seu aplicativo.**

Para seus testes automatizados, isso significa que você não precisa apenas executar seu próprio aplicativo, mas também o componente com o qual está integrando. 
Se você estiver testando a integração com um banco de dados, precisará executar um banco de dados ao executar seus testes. 
Para testar se você pode ler arquivos de um disco, você precisa salvar um arquivo em seu disco e carregá-lo em seu teste de integração.

Um teste de integração de banco de dados ficaria assim:

1. iniciar um banco de dados
2. conecte seu aplicativo ao banco de dados 
3. acione uma função dentro do seu código que grava dados no banco de dados 
4. verifique se os dados esperados foram gravados no banco de dados lendo os dados do banco de dados

Outro exemplo, testar se seu serviço se integra a um serviço separado por meio de uma API REST pode ser assim:

1. inicie seu aplicativo
2. inicie uma instância do serviço separado (ou um teste duplo com a mesma interface)
3. acione uma função em seu código que lê a API do serviço separado
4. verifique se seu aplicativo pode analisar a resposta corretamente

Escreva testes de integração para todos os trechos de código em que você serializa ou desserializa dados. Exemplos:

- Chamadas para a API REST dos seus serviços
- Leitura e gravação em bancos de dados
- Chamada de APIs de outros aplicativos
- Leitura e gravação em filas
- Escrevendo no sistema de arquivos

Ao escrever testes de integração, você deve tentar executar suas dependências externas localmente: 
execute um banco de dados MySQL local, teste em um sistema de arquivos local, etc. 
Se você estiver integrando com um serviço separado, execute uma instância desse serviço localmente ou crie e execute uma versão falsa que imita o comportamento do serviço real.


#### Exemplo de teste de integração

[person/mysql/mysql_test.go](https://github.com/eminetto/post-testes-go/blob/main/person/mysql/mysql_test.go)

Este teste faz a validação da camada de integração com o banco de dados. 
Ele [cria um container Docker](https://github.com/eminetto/post-testes-go/blob/main/person/mysql/mysql_test.go#L18), 
[conecta no banco de dados](https://github.com/eminetto/post-testes-go/blob/main/person/mysql/mysql_test.go#L23), 
[cria as tabelas](https://github.com/eminetto/post-testes-go/blob/main/person/mysql/mysql_test.go#L28),
[executa os testes](https://github.com/eminetto/post-testes-go/blob/main/person/mysql/mysql_test.go#L36),
e no final [faz o truncate das tabelas](https://github.com/eminetto/post-testes-go/blob/main/person/mysql/mysql_test.go#L32) e [destrói o container](https://github.com/eminetto/post-testes-go/blob/main/person/mysql/mysql_test.go#L22)


#### Executando os testes de integração

Execute

    make integration


#### Suite test

Para reutilizar código e organizar a inicialização dos testes, pode ser utilizado o [Suite](https://github.com/stretchr/testify#suite-package) da própria Testify. Com o suite, podemos utilizar os métodos SetupTest e TearDownTest, garantindo um teste limpo e assertivo.

#### Exemplo de teste de integração usando o Suite

Neste [PR](https://github.com/eminetto/post-testes-go/pull/2) é possível ver o uso do Suite em uma versão do exemplo anterior.

### Teste end to end

Testes de ponta a ponta dão a você a maior confiança quando você precisa decidir se seu software está funcionando ou não. Mas devido ao alto custo de manutenção, você deve reduzir ao mínimo o número de testes completos.
Pense nas interações de alto valor que os usuários terão com seu aplicativo. Tente criar jornadas do usuário que definam o valor central do seu produto e traduza as etapas mais importantes dessas jornadas do usuário em testes automatizados de ponta a ponta.

#### Exemplos de teste end to end

[internal/http/echo/handler_e2e_test.go](https://github.com/eminetto/post-testes-go/blob/main/internal/http/echo/handler_e2e_test.go)

Este teste implementa o fluxo de cadastro e leitura de um usuário. 


#### Executando os testes de integração

Execute

    make e2e

## Testes na correção de bugs

Testes, especialmente os unitários, são ótimas ferramentas para usarmos no momento da correção de um bug. Idealmente, quando um erro é reportado um bom fluxo para se seguir é:

1. Escreva um cenário de testes que produza o erro
2. Resolva o problema no código fonte
3. Execute os testes para garantir que nenhum efeito colateral foi adicionado
4. Refatore o código fonte caso necessário
5. Execute os testes novamente e faça o deploy da nova versão.


## Evite a duplicação de testes

Agora que você sabe que deve escrever diferentes tipos de testes, há mais uma armadilha a ser evitada: duplicar testes em todas as diferentes camadas da pirâmide. 
Embora seu pressentimento possa dizer que não existem "muitos testes", isso não é uma verdade. Cada teste em seu conjunto de testes é bagagem adicional e não vem de graça. Escrever e manter testes leva tempo. Ler e entender o teste de outras pessoas leva tempo. E, claro, executar testes leva tempo.

Assim como no código de produção, você deve buscar a simplicidade e evitar a duplicação. No contexto da implementação de sua pirâmide de teste, você deve manter duas regras em mente:

1. Se um teste de nível superior detectar um erro e não houver falha no teste de nível inferior, você precisará escrever um teste de nível inferior
2. Empurre seus testes o mais baixo possível na pirâmide de testes

A primeira regra é importante porque os testes de nível inferior permitem restringir melhor os erros e replicá-los de maneira isolada. Eles serão executados mais rapidamente e ficarão menos inchados quando você estiver depurando o problema em questão. 

A segunda regra é importante para manter seu conjunto de testes rápido. Se você testou todas as condições com confiança em um teste de nível inferior, não há necessidade de manter um teste de nível superior em seu conjunto de testes. Ter testes redundantes se tornará irritante em seu trabalho diário pois o conjunto de testes será mais lento e você precisará alterar mais lugares quando alterar o comportamento do seu código.


## Escrevendo código de teste limpo

Assim como na escrita de código em geral, criar um código de teste bom e limpo exige muito cuidado. Aqui estão mais algumas dicas para criar um código de teste sustentável:

- O código de teste é tão importante quanto o código de produção. Dê-lhe o mesmo nível de cuidado e atenção. *"este é apenas um código de teste"* não é uma desculpa válida para justificar um código desleixado
- Teste uma condição por teste. Isso ajuda você a manter seus testes curtos e fáceis de raciocinar. Em Go podemos usar a construção `t.Run`, como [neste exemplo](https://github.com/eminetto/post-testes-go/blob/main/person/service_test.go#L16).
- Usar uma [estrutura bem definida](#estrutura-dos-testes) facilita a construção de testes limpos.
- A legibilidade importa. Não tente ser excessivamente DRY. A duplicação é aceitável, se melhorar a legibilidade. Tente encontrar um equilíbrio entre o código [DRY e DAMP](https://stackoverflow.com/questions/6453235/what-does-damp-not-dry-mean-when-talking-about-unit-tests?answertab=trending#tab-top)

## Referências

- [The Practical Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html)
- [Testes de unidade e TDD](https://www.slideshare.net/lcobucci/testes-de-unidade-e-tdd-solisc-2011)
- [Mocks Aren't Stubs](https://martinfowler.com/articles/mocksArentStubs.html)
- [“Don’t Mock What You Don’t Own” in 5 Minutes](https://hynek.me/articles/what-to-mock-in-5-mins/)
