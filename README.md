# Quake Log Parser

## Descrição

Este projeto é um parser para os logs do jogo Quake 3 Arena. Ele agrupa informações de cada partida, incluindo mortes, jogadores e modos de morte, e gera relatórios detalhados em formato JSON.

## Objetivo

O objetivo deste projeto é processar os logs do servidor Quake 3 Arena e extrair informações úteis sobre as partidas, como o número total de mortes, ranking de jogadores por número de kills, e as causas das mortes.

## Requisitos

- [Go](https://golang.org/doc/install) (versão 1.16 ou superior)
- Sistema operacional compatível com Go (Linux, macOS, Windows)

## Estrutura do Projeto

```bash
parser-golang/
├── main.go               # Ponto de entrada do programa
├── parser/               # Contém a lógica de parsing
│   ├── parser.go         # Parsing dos logs do Quake
│   ├── models.go         # Definições de estruturas de dados
│   └── constants.go      # Constantes usadas no projeto
├── logs/
│   └── quake_game.log    # Arquivo de log do jogo Quake 3 Arena (coloque aqui)
├── README.md             # Documentação do projeto
└── go.mod                # Gerenciamento de módulos do Go
```

## Como Executar

### 1. Clonar o Repositório

Clone o repositório para sua máquina local:

```bash
git clone https://github.com/seu-usuario/quake-log-parser.git
cd quake-log-parser
```

### 2. Preparar o Arquivo de Log

Coloque o arquivo de log que deseja analisar no diretório `logs/` com o nome `quake_game.log`. Certifique-se de que o arquivo esteja no formato adequado.

### 3. Executar o Script

Execute o script usando o comando:

```bash
go run main.go
```

Isso processará o arquivo de log e gerará um arquivo JSON chamado `quake_report.json` com os resultados.

## Exemplo de Uso

Após executar o script, você verá a saída em formato JSON com informações sobre cada partida. Exemplo de uma saída simplificada:

```json
{
  "game_1": {
    "total_kills": 45,
    "players": ["Dono da Bola", "Isgalamido", "Zeh"],
    "kills": {
      "Dono da Bola": 5,
      "Isgalamido": 18,
      "Zeh": 20
    },
    "kills_by_means": {
      "MOD_SHOTGUN": 10,
      "MOD_RAILGUN": 2,
      "MOD_GAUNTLET": 1
    }
  }
}
```

## Funcionalidades Implementadas

1. **Parsing de Logs**: Lê e analisa os logs do Quake 3 Arena.
2. **Agrupamento por Partidas**: Agrupa os dados por partida e extrai informações como total de mortes, jogadores e kills.
3. **Relatório de Mortes por Modos**: Gera um relatório de mortes agrupadas por causa (modos de morte).
4. **Ranking de Jogadores**: Cria um ranking dos jogadores com base no número de kills.

## Considerações Finais

Este projeto foi desenvolvido como parte de um teste para a posição de Security Analyst|Specialist|Engineer na CloudWalk. Ele visa demonstrar a capacidade de trabalhar com parsing de logs e gerar relatórios detalhados a partir dos dados extraídos.

Para dúvidas ou sugestões, entre em contato através do hevertydouradob@gmail.com ou abra uma issue no repositório.

## Licença

None.
```
