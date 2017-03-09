package providers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"unicode"

	"bytes"

	"github.com/PuerkitoBio/goquery"
)

const mangaHostURL = "http://mangahost.net/manga"

/*Note to self:
    A good design is to make your type unexported, but provide an exported constructor function like NewMyType()
in which you can properly initialize your struct / type. Also return an interface type and not a concrete type,
and the interface should contain everything others want to do with your value. And your concrete type must
implement that interface of course.
*/

// UnionMangas é uma estrutura pra representar o download de paginas provenientes da union mangas
type mangaHost struct {
	nomeMangaFormatada string
}

func (u *mangaHost) errNotFound() {
	log.Println("Não foi possivel encontrar o mangá. Determinando o problema.")
	bc := Capitulo
	Capitulo = "1"
	resh, err := goquery.NewDocument(u.GerarURL())
	handle(err)
	if resh.Find(".error404").Length() > 0 {
		log.Fatal("Mangá ", MangaAtual, " não existe")
	} else {
		log.Fatal("Capitulo ", bc, " não existe")
	}
}

func (u *mangaHost) formatarNome() {
	if u.nomeMangaFormatada != MangaAtual {
		new := ""
		val := strings.Replace(MangaAtual, " ", "-", 10)
		for _, st := range val {
			if unicode.IsLetter(st) || unicode.IsNumber(st) || string(st) == "-" {
				new += string(st)
			}
		}
		u.nomeMangaFormatada = new
		MangaAtual = new
	}
}

//GerarURL retorna a url da union Mangas
func (u *mangaHost) GerarURL() string {
	u.formatarNome()
	capNum, err := strconv.Atoi(Capitulo)
	if err != nil || capNum <= 0 {
		log.Fatalf("%s não é um numero de capitulo válido", Capitulo)
	}
	return fmt.Sprintf("%s/%s/%d", mangaHostURL, u.nomeMangaFormatada, capNum)
}

func (u *mangaHost) ListImgURL() []string {
	resh, err := goquery.NewDocument(u.GerarURL())
	handle(err)
	if resh.Find(".error404").Length() > 0 {
		u.errNotFound()
	}
	urls := make([]string, 0)
	html, err := resh.Html()

	lines := strings.Split(html, "\n")
	bullseye := "var images = [\"<a href='"
	for ln := len(lines) - 1; ln >= 0; ln-- {
		line := strings.TrimSpace(lines[ln])
		if len(line) < len(bullseye) {
			continue
		}
		if line[:len(bullseye)] == bullseye {
			bullseye = line[15 : len(line)-3]
			bullseye = strings.Replace(bullseye, "\",\"", " \n", -1)
			break
		}
	}
	resh, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(bullseye)))
	handle(err)
	paginas := resh.Find("img[id]")
	log.Println(paginas.Length(), "páginas encontradas, disponiveis, possiveis duplicadas. Filtrando.")
	uniqPages := make(map[string]string)
	paginas = paginas.Each(func(a int, sel *goquery.Selection) {
		imgURL, exis := sel.Attr("src")
		if !exis {
			log.Panic("Mudança no layout. propriedadade src nao contem mais a URL no MangaHost. Mande um commit com isso")
			return
		}
		// Filtra a URL pra pegar o nome do arquivo sem qualquer extenção
		noExtName := imgURL[strings.LastIndex(imgURL, "/")+1:]
		noExtName = noExtName[:strings.Index(noExtName, ".")]
		// Se o arquivo nao está no mapa ou é um PNG
		if uniqPages[noExtName] == "" {
			uniqPages[noExtName] = imgURL
		} else if strings.Contains(imgURL, noExtName+".png") {
			old := uniqPages[noExtName]
			uniqPages[noExtName] = imgURL
			// Se nao for um numero já adiciona diretamente
			_, err := strconv.Atoi(noExtName)
			if err != nil {
				urls = append(urls, imgURL)
				for l := 0; l < len(imgURL); l++ {
					if urls[l] == old {
						urls = append(urls[:l], urls[l+1:]...)
					}
				}
			}
		}
	})
	for k := 1; k <= len(uniqPages); k++ {
		if uniqPages[fmt.Sprintf("%02d", k)] != "" {
			urls = append(urls, uniqPages[fmt.Sprintf("%02d", k)])
		}
	}
	if paginas.Length() != len(urls) {
		log.Printf("%d duplicadas encontradas e filtradas. %d paginas totais. \n",
			(paginas.Length() - len(uniqPages)), len(urls))
	}
	return urls
}

//TtlCapitulos retorna o total de capitulos disponiveis
func (u *mangaHost) TTLCapitulos() string {
	bak := Capitulo
	Capitulo = "1"
	doc, err := goquery.NewDocument(u.GerarURL())
	handle(err)
	last, exist := doc.Find(".viewerChapter").First().First().Attr("value")
	println(last)
	if !exist {
		log.Fatal("Mangá especificado não existe.")
	}

	Capitulo = bak
	return last
}

func (u *mangaHost) PesquisarTitulos(manga string) []string {
	// Constroi a url e a variavel pra receber o valor final
	resul, err := goquery.NewDocument("http://mangahost.net/find/" + manga)
	handle(err)

	temp := make([]string, 0)
	possib := resul.Find("tbody").Find("a.pull-left")
	possib.Each(func(i int, r *goquery.Selection) {
		m, exts := r.Attr("title")
		if !exts {
			log.Fatalln("Engine de pesquisa manga host foi alterada.")
		}
		temp = append(temp, m)
	})

	return temp
}

func (u *mangaHost) Nome() string {
	return "MangáHost"
}
