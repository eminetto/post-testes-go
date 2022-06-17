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

`GET /hello: Retorna "Hello World!". `
`GET /hello/{lastname}: Procura no banco de dados a pessoa pelo seu sobrenome e retorna "Hello {Firstname} {Lastname}" se a pessoa é encontrada. Retorna 404 caso não encontrada.`
`GET /weather: Chama uma API de previsão do tempo via HTTP e retorna as condições de Florianópolis, Brasil` ;)


## Arquitetura interna

Para este exemplo está sendo usada a [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) e o código é baseado no apresentado [neste post](https://eltonminetto.dev/post/2020-06-29-clean-architecture-2anos-depois/) e neste [repositório](https://github.com/eminetto/clean-architecture-go-v2)