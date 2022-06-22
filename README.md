# Workshop sobre testes

Repositório criado para exemplificar os tipos de testes e boas práticas que podem ser aplicados.

Baseado [neste post](https://martinfowler.com/articles/practical-test-pyramid.html) e na [aplicação](https://github.com/hamvocke/spring-testing) usada como exemplo.


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

A aplicação fornece três endponts:

```

GET /hello: Retorna "Hello World!". 
GET /hello/{lastname}: Procura no banco de dados a pessoa pelo seu sobrenome e retorna "Hello {Firstname} {Lastname}" se a pessoa é encontrada. Retorna 404 caso não encontrada.
GET /weather/{lat}/{long}: Chama uma API de previsão do tempo via HTTP e retorna as condições de acordo com as coordenadas. Retorna 404 caso não encontrada.

```


## Arquitetura interna

Para este exemplo está sendo usada a [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) e o código é baseado no apresentado [neste post](https://eltonminetto.dev/post/2020-06-29-clean-architecture-2anos-depois/) e neste [repositório](https://github.com/eminetto/clean-architecture-go-v2)

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

Antes de mergulhar nos tipos de teste,  uma boa estrutura para todos os testes é esta:

1. Configure os dados de teste
2. Invoque o método/função sendo testada
3. Confirme que os resultados esperados são retornados

Vamos observar esta estrutura em todos os testes.

### Testes unitários

Testes de unidade garantem que uma determinada unidade (o *sujeito em teste*) da base de código funcione conforme o esperado. Os testes de unidade têm o escopo mais restrito de todos os testes do conjunto de testes. O número de testes de unidade do conjunto de testes superará em grande parte qualquer outro tipo de teste.

#### Exemplos de teste unitário

[usecase/person/service_test.go](https://github.com/PicPay/go-test-workshop/blob/main/usecase/person/service_test.go)

Este arquivo contém os testes do serviço que implementa a interface [UseCase](https://github.com/PicPay/go-test-workshop/blob/main/usecase/person/interface.go#L33). 

Como o serviço tem por dependência uma implementação da interface [Repository](https://github.com/PicPay/go-test-workshop/blob/main/usecase/person/interface.go#L17) (que por sua vez precisa de uma conexão com o banco de dados), vamos usar o conceito de [mocks](https://martinfowler.com/articles/mocksArentStubs.html) para mantermos o foco do teste apenas na regra de negócio do serviço.
Para gerarmos facilmente os `mocks` estamos usando a ferramenta [mockery](https://github.com/vektra/mockery), que lê as interfaces e gera código para usarmos nos testes.
A geração dos `mocks` é executada pelo comando `make generate-mocks` e pode ser executada manualmente ou automaticamente quando executamos o comando `make unit-test`


[usecase/weather/service_test.go](https://github.com/PicPay/go-test-workshop/blob/main/usecase/weather/service_test.go)

Este arquivo contém os testes do serviço que implementa a interface [UseCase](https://github.com/PicPay/go-test-workshop/blob/main/usecase/weather/interface.go#L12).

Este é um serviço que faz uso de uma [API externa](https://api.openweathermap.org/). Para não acessar a API real a cada teste criamos um `mock` para simular o seu comportamento. 
Vale destacar uma boa prática neste pacote. Ao invés de colocarmos como dependência do `UseCase` um `http.Client` padrão da linguagem foi criada uma [interface](https://github.com/PicPay/go-test-workshop/blob/main/usecase/weather/interface.go#L8) para ser usada como dependência. 
No [construtor](https://github.com/PicPay/go-test-workshop/blob/main/usecase/weather/service.go#L20) do serviço criamos uma instância de `http.Client` e damos a opção do usuário substituir esse cliente padrão por outra implementação. 
Fazemos uso desta opção no [momento do teste](https://github.com/PicPay/go-test-workshop/blob/main/usecase/weather/service_test.go#L29) ao passar um `mock` do client. 
Esta implementação pode ser resumida pela frase `“Don’t Mock What You Don’t Own”` e mais detalhes podem ser vistos neste [post](https://hynek.me/articles/what-to-mock-in-5-mins/).


[api/handler_test.go](https://github.com/PicPay/go-test-workshop/blob/main/api/handler_test.go)

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

[infraestructure/repository/person/mysql_test.go](https://github.com/PicPay/go-test-workshop/blob/main/infraestructure/repository/person/mysql_test.go)

Este teste faz a validação da camada de integração com o banco de dados. 
Ele [cria um container Docker](https://github.com/PicPay/go-test-workshop/blob/main/infraestructure/repository/person/mysql_test.go#L18), 
[conecta no banco de dados](https://github.com/PicPay/go-test-workshop/blob/main/infraestructure/repository/person/mysql_test.go#L23), 
[cria as tabelas](https://github.com/PicPay/go-test-workshop/blob/main/infraestructure/repository/person/mysql_test.go#L28),
[executa os testes](https://github.com/PicPay/go-test-workshop/blob/main/infraestructure/repository/person/mysql_test.go#L36),
e no final [faz o truncate das tabelas](https://github.com/PicPay/go-test-workshop/blob/main/infraestructure/repository/person/mysql_test.go#L32) e [destrói o container](https://github.com/PicPay/go-test-workshop/blob/main/infraestructure/repository/person/mysql_test.go#L22)


#### Executando os testes de integração

Execute

    make integration


### Teste end to end

Testes de ponta a ponta dão a você a maior confiança quando você precisa decidir se seu software está funcionando ou não. Mas devido ao alto custo de manutenção, você deve reduzir ao mínimo o número de testes completos.
Pense nas interações de alto valor que os usuários terão com seu aplicativo. Tente criar jornadas do usuário que definam o valor central do seu produto e traduza as etapas mais importantes dessas jornadas do usuário em testes automatizados de ponta a ponta.

#### Exemplos de teste end to end

[api/handler_e2e_test.go](https://github.com/PicPay/go-test-workshop/blob/main/api/handler_e2e_test.go)

Este teste implementa o fluxo de cadastro e leitura de um usuário. 


#### Executando os testes de integração

Execute

    make e2e



## Referências

[The Practical Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html)