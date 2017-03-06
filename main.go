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

	"image/jpeg"

	"github.com/samosaara/gomanga/providers"
	"golang.org/x/image/webp"
)

//Provedor é a URL
var Provedor providers.Provedor

//Pasta é a localização dos downloads
var Pasta = "./gomanga/"

// Substituir é a que controla se arquivos existentes devem ser substituidos
var Substituir = false

// Espacos dita se os caminhos das imagems devem conter espaços
var Espacos = false

func main() {
	handleArgs()
	defaults()
	log.Printf("Agradecimento a %s por hostear o mangá que você esta prestes a baixar.", Provedor.Nome())
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
	capitulo := flag.String("c", "0", "Numero do `capitulo`. 0 para o mais novo.")
	prov := flag.String("site", "def", "Especifica o `site` de onde baixar/pesquisar. UnionMangas é o padrão. site pode ser:"+
		"\n\t- br_mh para MangaHost"+
		"\n\t- br_um para UnionMangas")
	help := flag.Bool("h", false, "Mostra essa ajuda e sai(exclusivo)")
	pesquisa := flag.String("p", "", "Pesiquisa `mangá` mostra resultados e sai(exclusivo)")

	flag.Parse()
	if flag.NArg() > 0 {
		log.Fatal("Argumento perdido sem operador.")
	}

	switch *prov {
	case "br_um", "def":
		Provedor = providers.UnionMangas
	case "br_mh":
		Provedor = providers.MangaHost
	default:
		log.Fatalln("Site provedor especificado inválido")
	}

	if *pesquisa != "" && !*help {
		if *prov == "def" {
			search(*pesquisa, nil)
		} else {
			search(*pesquisa, Provedor)
		}
		os.Exit(0)
	}

	if *help || providers.MangaAtual == "" {
		fmt.Println("Uso: gomanga [-s] [-r] [-c <num_manga>] -m \"Nome Manga\" OU",
			"gomanga -p \"pesquisa\"")
		flag.PrintDefaults()
		fmt.Println("Operadores exclusivos, são usados sozinhos e inibem o funcionamento dos outros.")
		os.Exit(0)
	}

	capt, err := strconv.Atoi(*capitulo)
	if err != nil && strings.Contains(*capitulo, "-") {
		providers.Capitulo = *capitulo
	} else if capt < 0 {
		log.Fatal("Numero de capitulo invalido")
	} else if capt == 0 {
		providers.Capitulo = ""
	} else {
		providers.Capitulo = *capitulo
	}
}

func search(manga string, src providers.Provedor) {
	mangas := make(map[string][]string)
	log.Println("Enviando pesquisa ao servidor.")

	if src == nil {
		log.Printf("Provedor nao especificado, pesquisando em todos")
		for _, prov := range providers.Provedores {
			name := prov.Nome()
			log.Printf("Pesquisando no site %s", name)
			mangas[name] = prov.PesquisarTitulos(manga)
		}
	} else {
		mangas[src.Nome()] = src.PesquisarTitulos(manga)
	}
	for provid, results := range mangas {
		if len(results) > 0 {
			fmt.Println("Mangas disponiveis encontrados em", provid, "foram:")
			for i, v := range results {
				fmt.Printf("  - %02d) %s\n", i+1, v)
			}
			fmt.Printf("Especifique esse provedor com a opcao -site")
		} else {
			fmt.Printf("Nenhum mangá encontrado em %s\n", provid)
		}
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

func createIfNotExists(file string) (*os.File, bool) {
	r, err := os.Open(file)
	// create file if not exists
	if os.IsNotExist(err) || Substituir {
		r, err = os.Create(file)
		handle(err)
		return r, true
	} else if err == nil {
		return r, false
	}
	handle(err)
	return nil, false
}

// download é o que parece ser. Baixa o mangá, dada a lista de URLs proveniente do provedor
func download() {
	log.Println("Aguarde, baixando lista de páginas do capitulo", providers.Capitulo, "para download.")
	imgs := Provedor.ListImgURL()

	local := Pasta + "/" + Provedor.Nome() + "/"
	if Espacos {
		local += strings.Replace(providers.MangaAtual, "_", " ", -1)
	} else {
		local += strings.Replace(providers.MangaAtual, " ", "_", -1)
	}
	local += "/Capitulo-" + providers.Capitulo + "/"
	err := os.MkdirAll(local, 0777)
	handle(err)
	log.Printf("Econtradas %d paginas disponiveis, iniciando download da primeira.", len(imgs))

	for i, imgURL := range imgs {
		originalName := imgURL[strings.LastIndex(imgURL, "/")+1:]
		originalExt := originalName[strings.LastIndex(originalName, "."):]

		destFilename := fmt.Sprintf("Pagina_%02d", i+1)
		destPath := local + destFilename

		var dst *os.File
		var isNew bool

		//Novos edge-cases de extensões diferenciadas podem ser adicionadas sem alterar muito do código
		if originalExt == ".webp" {
			destFilename += ".jpg"
			fFile, isNew := createIfNotExists(destPath + ".jpg")

			if isNew {
				img, err := http.Get(imgURL)
				handle(err)
				log.Println("Página está em .webp baixando e convertendo antes de escrever ao disco.")
				webpImg, err := webp.Decode(img.Body)
				handle(err)
				jpeg.Encode(fFile, webpImg, nil)
				log.Printf("%s convertida para JPG e escrita no disco com sucesso. Iniciando download da prox. pagina.", destFilename)
				img.Body.Close()
				fFile.Close()
				continue
			}
		} else {
			dest, isNew := createIfNotExists(destPath + originalExt)
			destFilename += originalExt

			if isNew {
				dst = dest
			}
		}

		if dst != nil {
			img, err := http.Get(imgURL)
			handle(err)
			_, err = io.Copy(dst, img.Body)
			handle(err)
			log.Printf("%s completa, iniciando o download da próxima pag.", destFilename)
			dst.Close()
			img.Body.Close()
			continue
		}
		if !isNew && !Substituir {
			log.Printf("Arquivo %s já existe. Ignorando. Adicione -r para substituir.\n", destFilename)
		}
	}
}
