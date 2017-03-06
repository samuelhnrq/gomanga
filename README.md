# GoManga
Cliente para **download** de mangás provenientes dos sites: [Union Mangas](http://unionmangas.net) e [Mangá Host](http://mangahost.net/). Por enquanto é bem bem simples e só prova de conceito, por mais já seja completamente funcional. Esse projeto é matido por mim como hobby e pet-project. Qualquer pull-request e/ou issue é  muito bem vindo.

## Instalação
```
go get github.com/samosaara/gomanga
```

## Usagem
Existem os argumentos de download e os argumentos de uso exclusivo que devem ser usados sozinhos.

### Download
``` 
gomanga [-c numero_do_capitulo] [-r] [-s] [-site provedor] -m Nome_Manga
```
- Nome_manga é o mesmo da url do site da union mangas pro mangá desejado. Exemplo, Dragon Ball capitulo 20 seria: `gomanga  -c 20 -m Dragon_Ball`
- Os mangás são sempre salvos em `./gomangas/{Nome_manga}/Capitulo{NumCapitulo}/`
- A maioria dos nomes é bem obvil. E não diferenciam maiuscula e minuscula. Coloque o nome do mangá entre aspas se tiver espaços.
- Ordem dos argumentos não faz diferença
- A opção 'site' Especifica o site de onde baixar o manga. Pode ser usada na pesquisa ou no download. Deve ser um desses:
  - 'br_um' para Union Mangás
  - 'br_mh' para Mangá Host
- A opção 'c' pode ser ou um numero ou um intervalo de capitulos. ex: `-c 9` ou `-c 5-9`
- A opção 'r' baixa e substitui mesmo se os arquivos já existirem.
- A opção 's' salva coloca espaços em vez de underlines no nome das pastas dos mangás

### Exclusivos
- A opção 'p' em conjunto com uma string pesquisa o mangá dos provedores disponiveis.
- Falta de opções ou a opção 'h' mostra a ajuda

## Inspiração e agradecimentos:
[kumroute/unionmangas](https://github.com/kumroute/unionmangas)
