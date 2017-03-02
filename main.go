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
var Pasta = "./Manga/"

// MangaAtual é o nome do manga dessa vez
var MangaAtual = ""

// Capitulo é o Capitulo atual
var Capitulo = ""

func main() {
	handleArgs()
	defaults()
	download(buildURL())
}

func handleArgs() {
	args := os.Args[1:]
	var opt, val string
	checkEmpty := func() {
		if val == "" {
			log.Fatal("Argumento ", opt, " exige um valor")
		}
	}
	for i := 0; i < len(args); i++ {
		if string(args[i][0]) != "-" {
			log.Fatal("Argumento perdido")
		}
		opt = string(args[i][1:])
		if string(args[i+1][0]) != "-" {
			val = args[i+1]
			args = args[2:]
		} else {
			val = ""
			args = args[1:]
		}
		switch opt {
		case "m", "manga":
			checkEmpty()
			strings.Replace(val, " ", "_", -1)
			MangaAtual = val
		case "c", "capitulo":
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
		}
		opt, val = "", ""
	}
}

func defaults() {
	if Provider == "" {
		Provider = UNIONMANGAS
	}
	if MangaAtual == "" {
		log.Fatal("Manga precisa ser especificado")
	}
	if Capitulo == "" {
		Capitulo = acharUltimoCap()
	}
}

func acharUltimoCap() string {
	doc, err := goquery.NewDocument(buildURL())
	handle(err)
	last, exist := doc.Find("#cap_manga1").Children().Last().Attr("value")
	if !exist {
		log.Fatal("Mudança por parte da union mangas. Terei que atualizar meu codigo")
	}

	return last
}

func buildURL() string {
	var url bytes.Buffer

	url.WriteString(Provider)
	url.WriteString("/")
	url.WriteString(MangaAtual)
	url.WriteString("/")
	url.WriteString(Capitulo)
	return url.String()
}

func handle(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func download(url string) {
	resh, err := http.Get(url)
	handle(err)
	defer resh.Body.Close()

	res, err := html.Parse(resh.Body)
	handle(err)
	paginas := acharKey(res, "img", "class", "img-responsive")
	pageNum := 1
	for _, pag := range paginas {
		var imgUrl string
		for _, atr := range pag.Attr {
			if atr.Key == "data-lazy" {
				imgUrl = atr.Val
				break
			}
		}
		if strings.Contains(imgUrl, "http://unionmangas.net/images") || strings.Contains(imgUrl, ".gif") {
			continue
		}

		img, err := http.Get(imgUrl)
		handle(err)
		defer img.Body.Close()

		local := Pasta + MangaAtual
		err = os.MkdirAll(local, 0777)
		handle(err)
		if pageNum < 10 {
			local += "/Pagina0"
		} else {
			local += "/Pagina"
		}
		local += strconv.Itoa(pageNum) + string(imgUrl[strings.LastIndex(imgUrl, "."):])
		pageNum++
		file, err := os.Create(local)
		handle(err)

		log.Printf("Baixando a pagina %d ", pageNum-1)

		_, err = io.Copy(file, img.Body)
		handle(err)
		file.Close()
	}
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
