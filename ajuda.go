package main

import "fmt"

func ajuda() {
	fmt.Println("Uso: gomanga [-s] [-r] [-c <num_manga>] -m \"Nome Manga\"\n" +
		"Baixa o mangá especificado numa nova pasta dentro da atual chamada gomanga, contendo" +
		"uma pasta com o nome do mangá que por sua vez contem uma pasta com todas as páginas de" +
		"cada capitulo\n\n" +
		"Argumentos Obrigatórios:\n" +
		"\t-m <manga> Especifica o mangá a ser baixado.\n" +
		"Argumentos Opicionais:\n" +
		"\t-s Adiciona espaços na pasta destino do mangá\n" +
		"\t-r Substitui se os arquivos já existirem.\n" +
		"\t-c <num> Especifica o capitulo a ser baixado, baixa o ultimo disponivel se omitido.\n")
}
