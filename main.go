package main

import (
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
			if first < 10 {
				providers.Capitulo = "0" + strconv.Itoa(first)
			} else {
				providers.Capitulo = strconv.Itoa(first)
			}
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
	args := os.Args[1:]
	var opt, val string
	if len(args) == 0 {
		ajuda()
		os.Exit(0)
	}
	checkEmpty := func() {
		if val == "" {
			log.Fatal("Argumento", opt, "exige um valor")
		}
	}
	precisaVazio := func() {
		if val != "" {
			log.Fatal("Argumento", opt, "não aceita valor algum deve ser usado sozinho.")
		}
	}
	for len(args) > 0 {
		if string(args[0][0]) != "-" {
			log.Fatal("Argumento perdido")
		}
		opt = string(args[0][1:])
		if len(args) > 1 {
			if string(args[1][0]) != "-" {
				val = args[1]
				args = args[2:]
			} else {
				val = ""
				args = args[1:]
			}
		} else {
			val = ""
			args = args[1:]
		}
		switch opt {
		case "m", "-manga":
			checkEmpty()
			providers.MangaAtual = val
		case "c", "-capitulo":
			checkEmpty()
			if strings.Contains(val, "-") {
				providers.Capitulo = val
				break
			}
			capNum := 0
			capNum, err := strconv.Atoi(val)
			if err != nil || capNum <= 0 {
				log.Fatal("Numero de Capitulo invalido")
			}
			if capNum < 10 {
				providers.Capitulo = "0" + strconv.Itoa(capNum)
			} else {
				providers.Capitulo = strconv.Itoa(capNum)
			}
		case "r", "-substituir":
			precisaVazio()
			Substituir = true
		case "s", "-espacos":
			precisaVazio()
			Espacos = true
		}
		opt, val = "", ""
	}
}

func defaults() {
	if providers.MangaAtual == "" {
		log.Println("Manga precisa ser especificado.")
		log.Println("Use '_' para espaços. Precisa ser o mesmo nome da URL do manga. Não diferencia maiscula ou minuscula.")
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
