package providers

import "log"

// provedor é uma interface para unificar os diferentes provedores
type provedor interface {
	GerarURL() string
	PesquisarTitulos(manga string) []string
	TTLCapitulos() string
	ListImgURL() []string
}

// Capitulo é o valor do capitulo
var Capitulo string

// MangaAtual é o nome do manga escolhido atual
var MangaAtual string

// UnionMangas Representa o union mangas
var UnionMangas provedor = &unionMangas{""}

// handle é uma função para facilitar tratar de erros indesejados e evitar verbosidade
func handle(erros ...error) {
	for _, err := range erros {
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
