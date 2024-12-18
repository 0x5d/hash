# Hash

Una prueba de concepto de un acortador de URLs.

Vea el [documento de diseño](DOC.md) para más detalle.

Estructura:
- `/src` contiene el código del servidor HTTP.
- `/tf` contiene los módulos de Terraform para desplegar la infraestructura de la prueba de concepto.

## Ejecutar localmente

### Prerrequisitos

- docker
- [ko](https://ko.build/)

```sh
./build.sh
docker compose up
```

El servidor web estará escuchando en el puerto 8080:
```sh
$ curl localhost:8080/url/ -d '{"url": "https://www.mercadolibre.com.co/", "enabled": true}'  
{"shortened":"http://localhost:8080/b","enabled":true}
```

Puede poner la URL devuelta en la respuesta en la barra de búsqueda del navegador, o usar `curl -L`:
```sh
$ curl -L http://localhost:8080/b 
<html>
  <head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <link rel="stylesheet" type="text/css" href="https://http2.mlstatic.com/ui/navigation/2.3.5/mercadolibre/navigation.css">
  <style>
    .ui-empty-state {
      position: relative;
      min-height: 25em
    }
```
Vea la documentación del [API HTTP](src/README.md) para más información.
