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
REDIS_URL=redis://redis:6379
```
Возможен запуск нескольких экземпляров back сервиса на разных портах. 
В этом случае нужно откорретировать конфиг nginx [nginx.conf](frontend%2Fnginx.conf)
что бы трафик балансировался на сервисы, а именно нужно добавить в `upstream backend_api` строку с портом на котором работают 
дополнительные экземпляры бэка. 
При запуске бэка в докере можно скейлить бэк ``docker-compose up --scale backend=2 -d``

### ⚠️ возможна ошибка при запуске докера
> npm error code EIDLETIMEOUT
> npm error Idle timeout reached for host registry.npmjs.org:443

ошибка связана с блокировкой ресурса https://registry.npmjs.org. Что бы обойти ошибку нужен ВПН или попробовать
расскоментировать строку `RUN npm config set registry https://registry.npmmirror.com` (установка зеркала)
в файле [Dockerfile](frontend%2FDockerfile).

### Архитектура 

```mermaid
flowchart LR
 subgraph subGraph0["Frontend (React)"]
        UI["Web UI (React)"]
  end
 subgraph subGraph1["Backend (Go)"]
        API["REST API & WebSocket"]
        Worker@{ label: "<pre style=\"font-family:'JetBrains\">Analyzer</pre>(Выполнение задач, запуск 1С через CLI)" }
        Cron["Scheduler"]
  end
 subgraph s1["1С"]
        OneC["1C:Enterprise CLI"]
  end
 subgraph Database["Database"]
        DB[("PostgreSQL")]
  end
 subgraph s2["Pub/Sub"]
        Redis[("Redis Pub/Sub")]
  end
    UI -- REST API --> API
    UI <-- WebSocket уведомления о прогрессе --> API
    API -- Команда на выгрузку --> Worker
    Cron -- Плановая выгрузка --> Worker
    Worker -- Запуск 1С --> OneC
    OneC -- Результаты выгрузки --> Worker
    API -- Чтение/Запись данных --> DB
    Worker -- Запись данных --> DB
    Worker -- Публикует прогресс --> Redis
    API -- Получает обновления прогресса (подписка) --> Redis
    Redis -- Push событий о прогрессе --> API

    Worker@{ shape: rect}
    style UI fill:#61dafb,stroke:#333,stroke-width:1px
    style API fill:#f0db4f,stroke:#333,stroke-width:1px
    style Worker fill:#ffe599,stroke:#333,stroke-width:1px
    style Cron fill:#fff2cc,stroke:#888,stroke-dasharray: 5 5
    style OneC fill:#efefef,stroke:#6fa8dc,stroke-width:1.5px
    style DB fill:#d9ead3,stroke:#38761d,stroke-width:1.5px
    style Redis fill:#cfe2f3,stroke:#b40101,stroke-width:2px




```