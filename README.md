# ttyslot

i wanted to make a terminal slot machine cuz I like to gamble but I don't have internet always

this is a 3 reel slot game built in go with some charm libraries for lots of beauty
![slot](https://github.com/trashdeveloper/ttyslot/blob/main/image.webp?raw=true)
## what it has

- 3 reel slot machine with weighted symbols
- smooth reel spin with staggered stopping
- payouts, pairs, cherries, and jackpot hits
- random free spin triggers
- bonus buy that scales with your current bet
- theme switching
- session stats and recent spin memory
- responsive layout so it still works when the terminal is not huge

## controls

- `space` / `enter` spin
- `up` / `down` change bet
- `b` buy bonus
- `t` switch theme
- `r` reset session
- `q` quit

## run it

```bash
go mod tidy
go run .
```

## dev stuff

```bash
make test
make build
make fmt
```

## terminal size

- nicest around `120x32`
- still works at `64x20`

## project layout

- `main.go` Bubble Tea app state, input handling, spin flow
- `game.go` symbol table, weighted rolls, payout logic
- `ui.go` layout, themes, styling, rendering
- `main_test.go` app/session tests
- `game_test.go` payout tests

## why i made it

I like go language because google made it and if u guys didn't know go lang starts with g o and so does google

## license

This project is licensed under GPLv3
