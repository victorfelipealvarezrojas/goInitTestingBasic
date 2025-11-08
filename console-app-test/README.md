
# Tipos que Funcionan por Referencia en Go

En Go, ciertos tipos funcionan de forma natural por referencia sin necesidad de marcarlos explícitamente con `*` (puntero) o `&` (dirección). Estos tipos incluyen:

## Channels
Los canales se pasan y manejan automáticamente como referencias. No requieren marcas de puntero para funcionar por referencia.

## Slices
Los slices son vistas sobre arrays subyacentes y ya se comportan por referencia de forma natural. Al pasarlos a funciones, se transmite la referencia al array subyacente.

## Maps
Los mapas se pasan por referencia automáticamente. No es necesario usar `*` o `&` para trabajar con ellos por referencia.

## Interfaces
Las interfaces contienen referencias internas y se pasan por referencia de forma automática sin necesidad de marcas explícitas.

## Funciones
Los valores de función ya son referencias y se comportan como tales naturalmente en Go.

## Nota
Para la mayoría de los otros tipos en Go, como structs, tipos primitivos y arrays, sí es necesario usar `*` y `&` si deseas pasar referencias explícitas.



# inicializar proyecto 
go mod init nombre_modulo  

# ejecutar test
go test .
go test -v .    

# crear archivo de covertura 
go test -coverprofile cover.out .

# visualizar en html 
go tool cover -html=cover.out    

