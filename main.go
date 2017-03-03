package main

import (
	"log"

	"net/http"

	"strings"

	"os"

	"io"

	"bytes"

	"strconv"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

//UNIONMANGAS é a URL
const UNIONMANGAS = "http://unionmangas.net/leitor"

//Provider é o URL atual
var Provider = ""

//Pasta é a localização dos downloads
var Pasta = "./gomanga/"

// MangaAtual é o nome do manga dessa vez
var MangaAtual = ""

// Capitulo é o Capitulo atual
var Capitulo = ""

// Substituir é a que controla se arquivos existentes devem ser substituidos
var Substituir = false

func main() {
	handleArgs()
	defaults()
	download(buildURL())
	log.Println("Obrigado por usar!")
}

func handleArgs() {
	args := os.Args[1:]
	var opt, val string
	checkEmpty := func() {
		if val == "" {
			log.Fatal("Argumento", opt, "exige um valor")
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
			//FIXME: Bem ineficiente, trocar por um split
			new := ""
			val = strings.Replace(val, " ", "_", 10)
			for _, st := range strings.Split(val, "_") {
				new += strings.Title(st) + "_"
			}
			new = string(new[:len(new)-1])
			val = new
			MangaAtual = val
		case "c", "-capitulo":
			checkEmpty()
			capNum, err := strconv.Atoi(val)
			if err != nil || capNum <= 0 {
				log.Fatal("Numero de Capitulo invalido")
			}
			if capNum < 10 {
				Capitulo = "0" + strconv.Itoa(capNum)
			} else {
				Capitulo = strconv.Itoa(capNum)
			}
		case "r", "-substituir":
			Substituir = true
		}
		opt, val = "", ""
	}
}

func defaults() {
	if Provider == "" {
		Provider = UNIONMANGAS
	}
	if MangaAtual == "" {
		log.Println("Manga precisa ser especificado.")
		log.Println("Use '_' para espaços. Precisa ser o mesmo nome da URL do manga. Não diferencia maiscula ou minuscula.")
		log.Fatal("ex: gomanga -m [nome_manga]")
	}
	if Capitulo == "" {
		Capitulo = acharUltimoCap()
	}
}

func acharUltimoCap() string {
	doc, err := goquery.NewDocument(buildURL())
	handle(err)
	log.Println("Capitulo não especificado, baixando lista para verificar o ultimo.")
	last, exist := doc.Find("#cap_manga1").Children().Last().Attr("value")
	if !exist {
		log.Fatal("Mudança por parte da union mangas. Terei que atualizar meu codigo")
	}
	log.Println("O ultimo capitulo disponivel é o " + last)

	return last
}

// Constroi a URL final, função presisaava ser refeita para aceitar varios
func buildURL() string {
	var url bytes.Buffer

	url.WriteString(Provider)
	url.WriteString("/")
	url.WriteString(MangaAtual)
	url.WriteString("/")
	url.WriteString(Capitulo)
	return url.String()
}

// handle é uma função para facilitar tratar de erros indesejados e evitar verbosidade
func handle(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// download é o que parece ser. Baixa o mangá, dada a URL, tinha que ser refeito com goquery
func download(url string) {
	log.Println("Baixando lista de páginas do capitulo", Capitulo, "para download.")
	resh, err := goquery.NewDocument(url)
	handle(err)

	paginas := resh.Find("img.img-responsive")
	log.Println(paginas.Length()-3, "páginas encontradas, disponiveis pra download")
	pageNum := 1
	paginas = paginas.Each(func(a int, sel *goquery.Selection) {
		imgURL, exis := sel.Attr("data-lazy")
		if !exis {
			println("URL de imagens não foi encontrada.")
			return
		}
		if strings.Contains(imgURL, "http://unionmangas.net/images") || strings.Contains(imgURL, ".gif") {
			return
		}

		local := Pasta + MangaAtual + "/Capitulo " + Capitulo
		err = os.MkdirAll(local, 0777)
		handle(err)
		local += string(imgURL[strings.LastIndex(imgURL, "/"):])
		pageNum++
		file, err := os.Open(local)
		if os.IsNotExist(err) || Substituir {
			file, err = os.Create(local)
		} else if err == nil {
			log.Println("Arquivo/Pagina", pageNum-1, "já existe... Pulando...")
			return
		}

		img, err := http.Get(imgURL)
		handle(err)
		defer img.Body.Close()

		log.Printf("Baixando a pagina %d ", pageNum-1)

		_, err = io.Copy(file, img.Body)
		handle(err)
		file.Close()
	})
}

// acharKey encontra e retorna os nodes que se encaixam com os valores
func acharKey(n *html.Node, tag string, prop string, valor string) []html.Node {
	finalSlice := make([]html.Node, 0)
	if n.Type == html.ElementNode && n.Data == tag {
		for _, a := range n.Attr {
			if a.Key == prop && (a.Val == valor || strings.Contains(a.Val, valor)) {
				finalSlice = append(finalSlice, *n)
				continue
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		finalSlice = append(finalSlice, acharKey(c, tag, prop, valor)...)
	}
	return finalSlice
}
