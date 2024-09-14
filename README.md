Utilizando GRPC com GOLANG

Utilizamos a lib https://github.com/ktr0731/evans?tab=readme-ov-file

Criamos um arquivo nome_do_arquivo.proto

Definimos os nossos inputs e outputs dos nossos seriços  gerar nossos proto
 
protoc --go_out=. --go-grpc_out=. path/nome_arquivo.proto

Depois disso criamos os nossos métodos no service

Depois de tudo definido podemos acessamos o nosso servidor

evans -r repl

# Mostrar nosso pacotes, temos que escolher um pacore para utilizar
show package

# Mostra nossos service 
show service

# Para executarmos uma chamada

service NomeDoService

Depois de selecionar o service, chamar o método

call NomeDoMetodo

