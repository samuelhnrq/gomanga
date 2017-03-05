package providers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/antonholmquist/jason"
)

const unionMangasURL = "http://unionmangas.net/leitor"

/*Note to self:
    A good design is to make your type unexported, but provide an exported constructor function like NewMyType()
in which you can properly initialize your struct / type. Also return an interface type and not a concrete type,
and the interface should contain everything others want to do with your value. And your concrete type must
implement that interface of course.
*/

// UnionMangas é uma estrutura pra representar o download de paginas provenientes da union mangas
type unionMangas struct {
	nomeMangaFormatada string
}

func (u *unionMangas) errNotFound() {
	log.Println("Não foi possivel encontrar o mangá. Determinando o problema.")
	bc := Capitulo
	Capitulo = "01"
	resh, err := goquery.NewDocument(u.GerarURL())
	handle(err)
	if resh.Url.Path == "/index.php" {
		log.Fatal("Mangá ", MangaAtual, " não existe")
	} else {
		log.Fatal("Capitulo ", bc, " não existe")
	}
}

func (u *unionMangas) formatarNome() {
	if u.nomeMangaFormatada != MangaAtual {
		new := ""
		val := strings.Replace(MangaAtual, " ", "_", 10)
		for _, st := range strings.Split(val, "_") {
			if st != "of" && st != "and" {
				new += strings.Title(st) + "_"
			} else {
				new += st + "_"
			}
		}
		new = string(new[:len(new)-1])
		u.nomeMangaFormatada = new
		MangaAtual = new
	}
}

//GerarURL retorna a url da union Mangas
func (u *unionMangas) GerarURL() string {
	u.formatarNome()
	capNum, err := strconv.Atoi(Capitulo)
	if err != nil {
		log.Fatalf("%s não é um numero de capitulo válido", Capitulo)
	}
	return fmt.Sprintf("%s/%s/%02d", unionMangasURL, u.nomeMangaFormatada, capNum)
}

func (u *unionMangas) ListImgURL() []string {
	resh, err := goquery.NewDocument(u.GerarURL())
	handle(err)
	if resh.Url.Path == "/index.php" {
		u.errNotFound()
	}
	urls := make([]string, 0)
	paginas := resh.Find("img.img-responsive")
	log.Println(paginas.Length()-3, "páginas encontradas, disponiveis pra download")
	paginas = paginas.Each(func(a int, sel *goquery.Selection) {
		imgURL, exis := sel.Attr("data-lazy")
		if !exis {
			log.Panic("Mudança no layout. propriedadade data-lazy nao contem mais a URL. Mande um commit com isso")
			return
		}
		if strings.Contains(imgURL, "http://unionmangas.net/images") || strings.Contains(imgURL, ".gif") {
			return
		}
		urls = append(urls, imgURL)
	})
	return urls
}

//TtlCapitulos retorna o total de capitulos disponiveis
func (u *unionMangas) TTLCapitulos() string {
	Capitulo = "01"
	doc, err := goquery.NewDocument(u.GerarURL())
	handle(err)
	last, exist := doc.Find("#cap_manga1").Children().Last().Attr("value")
	if !exist {
		log.Fatal("Mangá especificado não existe.")
	}

	return last
}

func (u *unionMangas) PesquisarTitulos(manga string) []string {
	// Constroi a url e a variavel pra receber o valor final
	jsonRes, err := http.Get("http://unionmangas.net/assets/busca.php?q=" + manga)
	handle(err)
	defer jsonRes.Body.Close()

	res, err := jason.NewObjectFromReader(jsonRes.Body)
	handle(err)

	temp := make([]string, 0)
	values, err := res.GetObjectArray("items")
	handle(err)

	for _, mangaObj := range values {
		mangaTitle, err := mangaObj.GetString("titulo")
		handle(err)
		temp = append(temp, mangaTitle)
	}

	return temp
}
