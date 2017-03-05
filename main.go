package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/samosaara/gomanga/providers"
)

//Provedor é a URL
var Provedor = providers.UnionMangas

//Pasta é a localização dos downloads
var Pasta = "./gomanga/"

// Substituir é a que controla se arquivos existentes devem ser substituidos
var Substituir = false

// Espacos dita se os caminhos das imagems devem conter espaços
var Espacos = false

func main() {
	handleArgs()
	defaults()
	if strings.Contains(providers.Capitulo, "-") {
		strRang := strings.Split(providers.Capitulo, "-")
		first, err := strconv.Atoi(strRang[0])
		last, err2 := strconv.Atoi(strRang[1])
		if err != nil || err2 != nil {
			log.Fatal("Intervalo de capitulos invalido")
		}
		log.Println("Conjunto de capitulos do", first, "ao", last, "especificados para dowload.")
		for ; first <= last; first++ {
			providers.Capitulo = fmt.Sprintf("%d", first)
			log.Println("Iniciando o download do capitulo", first)
			download()
			log.Println("Download do capitulo", first, "completo com sucesso")
		}
	} else {
		download()
	}
	log.Println("Obrigado por usar!")
}

func handleArgs() {
	flag.StringVar(&providers.MangaAtual, "m", "", "Especifica o nome do `mangá` a ser baixado")
	flag.BoolVar(&Substituir, "r", false, "Substitui arquivos existentes")
	flag.BoolVar(&Espacos, "s", false, "Coloca espaços no nome da pasta")
	capt := flag.Int("c", 0, "Numero do `capitulo`. 0 para o mais novo.")
	help := flag.Bool("h", false, "Mostra essa ajuda e sai(exclusivo)")
	pesquisa := flag.String("p", "", "Pesiquisa `mangá` mostra resultados e sai(exclusivo)")

	flag.Parse()
	if flag.NArg() > 0 {
		log.Fatal("Argumento perdido sem operador.")
	}

	if *pesquisa != "" && !*help {
		search(*pesquisa)
		os.Exit(0)
	}
	if *help || providers.MangaAtual == "" {
		fmt.Println("Uso: gomanga [-s] [-r] [-c <num_manga>] -m \"Nome Manga\" OU",
			"gomanga -p \"pesquisa\"")
		flag.PrintDefaults()
		fmt.Println("Operadores exclusivos, são usados sozinhos e inibem o funcionamento dos outros.")
		os.Exit(0)
	}
	if *capt < 0 {
		log.Fatal("Numero de capitulo invalido")
	} else if *capt == 0 {
		providers.Capitulo = ""
	} else {
		providers.Capitulo = strconv.Itoa(*capt)
	}
}

func search(manga string) {
	log.Println("Enviando pesquisa ao servidor.")
	mangas := Provedor.PesquisarTitulos(manga)
	if len(mangas) == 0 {
		log.Println("Nenhum mangá encontrado.")
		return
	}
	log.Println("Mangás encontrados:")
	for i, v := range mangas {
		fmt.Printf("  - %02d) %s\n", i+1, v)
	}
}

func defaults() {
	if providers.MangaAtual == "" {
		log.Println("Manga precisa ser especificado.")
		log.Println("Use '_' para espaços, ou coloque o nome do mangá entre aspas. Precisa ser o mesmo nome da URL do manga.",
			"Não diferencia maiscula ou minuscula.")
		log.Fatal("ex: gomanga -m [nome_manga]")
	}
	if providers.Capitulo == "" {
		log.Println("Capitulo não especificado, baixando lista para verificar o ultimo disponivel")
		providers.Capitulo = Provedor.TTLCapitulos()
		log.Println(providers.Capitulo, " é o ultimo capitulo disponivel")
	}
}

// handle é uma função para facilitar tratar de erros indesejados e evitar verbosidade
func handle(erros ...error) {
	for _, err := range erros {
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

// download é o que parece ser. Baixa o mangá, dada a lista de URLs proveniente do provedor
func download() {
	log.Println("Baixando lista de páginas do capitulo", providers.Capitulo, "para download.")
	imgs := Provedor.ListImgURL()
	local := Pasta
	if Espacos {
		local += strings.Replace(providers.MangaAtual, "_", " ", -1)
	} else {
		local += providers.MangaAtual
	}
	local += "/Capitulo-" + providers.Capitulo + "/"
	err := os.MkdirAll(local, 0777)
	handle(err)
	for i, imgURL := range imgs {
		imgName := "Pagina_"
		if i+1 < 10 {
			imgName += "0"
		}
		imgName += strconv.Itoa(i + 1)
		imgName += string(imgURL[strings.LastIndex(imgURL, "."):])
		imgFile := local + imgName
		file, err := os.Open(imgFile)
		if os.IsNotExist(err) || Substituir {
			file, err = os.Create(imgFile)
		} else if err == nil {
			log.Println("Arquivo", imgName, "já existe... Pulando...")
			continue
		}

		img, err := http.Get(imgURL)
		handle(err)

		log.Print("Baixando a ", imgName, "...")

		_, err = io.Copy(file, img.Body)
		handle(err)
		file.Close()
		img.Body.Close()

		log.Println("Download da", imgName, "completo com sucesso")
	}
}
