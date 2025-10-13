### extensions-info
Веб-приложение для разработчиков 1С, которое помогает анализировать расширения конфигурации.
Оно показывает, какие объекты метаданных переопределены в расширениях, в каких именно расширениях это происходит, и даже какие процедуры и функции были переопределены.

https://github.com/user-attachments/assets/3473fda7-108a-4188-a3f4-e927ee44c10f


### Запуск

Что бы полноценно все можно было запустить в докере нужно 1С 
затащить в докер, что б этого не делать можно все кроме бэка запустить в докере
командой
 ```
docker-compose up --build -d
 ```
и отдельно бэк 

* просто запустить (должен быть установлен [go](https://go.dev/dl/)) 
    ```
    make start-backend 
    ```
* собрать бинарный файл бэка (должен быть установлен [go](https://go.dev/dl/))
    ```
    make build-backend
    ```
* взять бинарник из [релиза](https://github.com/LazarenkoA/extensions-info/releases)

для запуска бэка нужны переменные окружения
```
PORT="8080"
POSTGRES_URL="postgres://postgres:password@localhost:5432/myapp?sslmode=disable"
```

### ⚠️ возможна ошибка при запуске докера
> npm error code EIDLETIMEOUT
> npm error Idle timeout reached for host registry.npmjs.org:443

ошибка связана с блокировкой ресурса https://registry.npmjs.org. Что бы обойти ошибку нужен ВПН или попробовать
расскоментировать строку `RUN npm config set registry https://registry.npmmirror.com` (установка зеркала)
в файле [Dockerfile](frontend%2FDockerfile).


#### TODO
- сделать поддержку возможности запуска нескольких экземпляров бэка. 
Сейчас приложение stateful из-за websocket. После поддержки stateles можно будет запускать несколько экземпляров
`docker-compose up --scale backend=2 -d`