package providers

import "log"

// Provedor é uma interface para unificar os diferentes provedores
type Provedor interface {
	GerarURL() string
	PesquisarTitulos(manga string) []string
	TTLCapitulos() string
	ListImgURL() []string
	Nome() string
}

// Capitulo é o valor do capitulo
var Capitulo string

// MangaAtual é o nome do manga escolhido atual
var MangaAtual string

// UnionMangas Representa o union mangas
var UnionMangas Provedor = &unionMangas{""}

// MangaHost É a instancia do provedor
var MangaHost Provedor = &mangaHost{""}

// Provedores é uma slice de todos provedores disponiveis
var Provedores = [...]Provedor{UnionMangas, MangaHost}

// handle é uma função para facilitar tratar de erros indesejados e evitar verbosidade
func handle(erros ...error) {
	for _, err := range erros {
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
