# Jantar dos Filósofos em Go

Quatro implementações do problema clássico do Jantar dos Filósofos, escritas em Go,
feitas para comparar na prática o que cada estratégia de sincronização custa em tempo
e em paralelismo real.

## O problema

Cinco filósofos sentam a uma mesa redonda. Entre cada par há **um** garfo — cinco
garfos no total. Um filósofo alterna entre pensar e comer, e para comer precisa dos
**dois** garfos ao seu lado. Como cada garfo é disputado por dois vizinhos, o problema
vira um estudo de exclusão mútua sob disputa circular.

A armadilha é a seguinte: se todos pegarem o garfo da esquerda ao mesmo tempo e ficarem
esperando o da direita, ninguém solta nada e o jantar trava para sempre. Esse é o
**deadlock**, e ele acontece porque as quatro condições de Coffman se satisfazem ao
mesmo tempo — exclusão mútua, posse e espera, não-preempção e espera circular. Quebrar
qualquer uma delas resolve, e é exatamente isso que as versões abaixo fazem de formas
diferentes.

## As quatro versões

| Pasta | Estratégia | Como quebra o deadlock |
|---|---|---|
| `Sequencial/` | **V0** — sem concorrência | Não há disputa: um filósofo come por vez. É a linha de base. |
| `ParaleloV1/` | **V1** — garfos como canais, sem proteção | **Não quebra.** É a versão ingênua, mantida de propósito para mostrar o problema. |
| `ParaleloV2/` | **V2** — mutex sobre a região crítica inteira | Elimina a espera circular serializando tudo: só um filósofo come por vez. |
| `ParaleloV3/` | **V3** — semáforo contador com N−1 fichas | Limita a 4 o número de filósofos disputando garfos. Com no máximo 4 competindo por 5 garfos, sempre há um que consegue os dois. |

Cada garfo é modelado como um `chan struct{}` de capacidade 1 contendo uma ficha:
receber do canal é pegar o garfo, enviar de volta é devolvê-lo. O semáforo da V3 usa a
mesma ideia, com um canal bufferizado de N−1 fichas.

## Resultados medidos

Uma execução de cada versão, 5 filósofos × 20 rodadas, 200 ms pensando e 500 ms comendo,
em Windows 11 com Go 1.27:

| Versão | Tempo total | Espera média pelo garfo | Maior espera | Comendo ao mesmo tempo |
|---|---|---|---|---|
| V0 — sequencial | 1m10,18s | — | — | 1 |
| V1 — sem proteção | 25,23s | 540,2 ms | 1,00s | 2 |
| V2 — mutex global | 50,27s | 1,762s | 2,00s | 1 |
| V3 — semáforo | 25,23s | 540,2 ms | 1,00s | 2 |

O que esses números dizem:

- **V2 é segura mas cara.** O mutex envolve pegar os garfos, comer e devolvê-los, então
  o paralelismo desaparece: 50s contra os 25s da V3, e a espera média triplica.
- **V3 entrega o melhor dos dois mundos.** Mesma velocidade da versão sem proteção, com
  a garantia formal de que não trava.
- **V1 e V3 empatam porque o deadlock da V1 quase nunca acontece.** A janela entre pegar
  um garfo e o outro é de nanossegundos, então a espera circular raramente se forma. Isso
  é justamente o que torna o bug perigoso: ele passa em todos os testes e aparece em
  produção.
- **V1 e V3 têm dois regimes.** Dependendo da ordem em que as goroutines acordam no
  primeiro ciclo, a execução pode cair num padrão de "comboio" em que só um filósofo come
  por vez — e o tempo vai para ~50s, igual ao da V2. O regime é decidido logo no início e
  se mantém até o fim. O semáforo evita o deadlock, mas não evita a fila.

## Verificação

As três versões paralelas foram executadas sob o detector de corridas do Go, sem nenhum
acesso concorrente não-sincronizado detectado:

```
$ go run -race ./ParaleloV1    # exit 0 — nenhum DATA RACE
$ go run -race ./ParaleloV2    # exit 0 — nenhum DATA RACE
$ go run -race ./ParaleloV3    # exit 0 — nenhum DATA RACE
```

Os contadores compartilhados (`numOfPhilosophersEating` e
`maxNumOfPhilosophersEatingSimult`) estão protegidos por mutex, e os garfos, por serem
canais, já carregam a sincronização necessária.

O repositório também passa limpo pelas ferramentas padrão:

```
$ gofmt -l .        # nenhuma saída
$ go vet ./...      # exit 0
$ go build ./...    # exit 0
```

No Windows, o `-race` exige cgo e um compilador C no PATH (por exemplo MinGW-w64);
sem isso o Go recusa a flag com `-race requires cgo`.

## Como rodar

Requer Go 1.21 ou superior (as versões paralelas usam o pacote `slices`).

```bash
go run ./Sequencial    # V0 — leva ~70s
go run ./ParaleloV1    # V1 — ~25s ou ~50s, conforme o regime
go run ./ParaleloV2    # V2 — ~50s
go run ./ParaleloV3    # V3 — ~25s ou ~50s, conforme o regime
```

Para compilar tudo de uma vez:

```bash
go build ./...
```

A saída mostra cada filósofo numa coluna, com `T<n>` para a rodada em que está pensando
e `E<n>` para a rodada em que está comendo, seguida das métricas agregadas ao final.

## Créditos

Trabalho da disciplina de Fundamentos de Processamento Paralelo e Distribuído (PUCRS),
desenvolvido em dupla com **[Felipe Torres](https://github.com/Felipe-Torres01)** — o
repositório original está em [Felipe-Torres01/T1FPPD](https://github.com/Felipe-Torres01/T1FPPD).

A estrutura base do laço pensar/comer e o formato de saída em colunas vêm do código de
exemplo do professor **César De Rose** (PUCRS).

## Licença

[MIT](LICENSE).
