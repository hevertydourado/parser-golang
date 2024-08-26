# Quake Log Parser

## Descrição

Este projeto é um parser para os logs do jogo Quake 3 Arena, agrupando informações de cada partida, incluindo mortes, jogadores e modos de morte.

## Como Executar

1. Clone o repositório.
2. Coloque o arquivo de log em `logs/quake_game.log`.
3. Execute o script:

```bash
go run main.go
```

## Estrutura do Projeto

```
quake-log-parser/
├── main.go
├── parser/
│   ├── parser.go
│   ├── models.go
│   └── constants.go
├── logs/
│   └── quake_game.log
├── README.md
└── go.mod
```
