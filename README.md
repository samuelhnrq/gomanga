# GoManga
Cliente para download de mangas proveniente do site unionmangas.net. Por enquanto é bem bem simples e só prova de conceito. Se tiver tempo adiciono mais funcionalidades. O pull-request é livre e muito bem vindo.

## Instalação
```
go get github.com/samosaara/gomanga
```

## Usagem
``` 
gomanga [-c numero_do_capitulo] [-r] -m Nome_Manga
```
- Nome_manga é o mesmo da url do site da union mangas pro mangá desejado. Exemplo, Dragon Ball capitulo 20 seria: `gomanga  -c 20 -m Dragon_Ball`
- A maioria dos nomes é bem obvil.
- A opção 'r' baixa e substitui mesmo se os arquivos já existirem.

#### Inspiração
[kumroute/unionmangas](https://github.com/kumroute/unionmangas)
