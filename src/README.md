# Server

El servidor HTTP contiene 2 endpoints descritos a continuación.

## Funcionamiento

Para crear una URL corta, se obtiene un ID monotónico de la base de datos y se traduce a base62, lo que se vuelve la "llave" de la URL. Luego, esta se almacena en la base de datos y la caché.

Al recibir una petición GET a `/<llave>`, primero se busca la caché. Si hay un hit, se retorna un estado 307 Moved Temporarily junto con la dirección objetivo registrada, lo cual le indica al cliente que debe redireccionar. Si no hay un hit, se decodifica la llave al ID numérico equivalente y se consulta en la base de datos.

## TODO

- Pruebas unitarias y de integración.

## HTTP API

### /

## GET
- Redirecciona (307) a la URL si esta existe, o
- Retorna Not Found (404) si no es así.
Ejemplo:
```sh
curl localhost:8080/b
```

### /url

#### POST /
Body:
```json
{"url": <URL>, "enabled": <bool>}
```
Respuesta:
```json
{"shortened": <shortened URL>, "enabled": <bool>}
```
Ejemplo:
```sh
curl localhost:8080/url/ -d '{"url": "https://www.mercadolibre.com.co/", "enabled": true}'
```

#### GET /<id>
Respuesta:
```json
{"shortened": <shortened URL>, "enabled": <bool>}
```
Ejemplo:
```sh
curl localhost:8080/url/b
```

#### PUT /<id>
Body:
```json
{"url": <URL>, "enabled": <bool>}
```
Respuesta:
```json
{"shortened": <shortened URL>, "enabled": <bool>}
```
Ejemplo:
```sh
 curl -X PUT localhost:8080/url/b -d '{"url": "mercadolibre.com"}'
```
