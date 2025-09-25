module github.com/Gitcelo/pokedexcli

go 1.24.4

require internal/pokeapi v1.0.0

replace internal/pokeapi => ./internal/pokeapi

require internal/pokecache v1.0.0

require (
	github.com/chzyer/readline v1.5.1 // indirect
	golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5 // indirect
)

replace internal/pokecache => ./internal/pokecache
