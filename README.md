# Оглавление

### Назначение сервиса:
Работа с метоположением пользователя на основе возможностей сервиса dadata
Основная возможность: поддкржка кеширования в redis + более удобный формат ответа

### Паспорт сервиса:
* go - rest-сервер на Golang
* redis - кеширование

#### Описание работы сервиса:
-

#### Настройка
* Поддерживаются 2 режима запуска: ```.docker/Dockerfile``` - prod, ```.docker/dev/Dockerfile``` - dev,
* В dev-режиме подключена библиотека github.com/githubnemo/CompileDaemon, позволяющая использовать hot reload code
* docker-compose up

#### Установка:
* создать ```.env``` по аналогии с ```.env.example```
* Проставить в файле .env значение CREDENTIALS_FROM_VAULT true, если хранение ключа на площадке предполагается в Vault и заполнить параметры подключения к Vault
* Запустить ```docker-compose up```

# DEVOPS

| Наименование        | Значение           | Примечание  |
| ------------- |:-------------| -----:|
| Name             | api                                               |                                      |
| Namespace        | geo-ns1-${branch}                       | ${branch} - master, developer, stage |
| Internal address | api.geo-ns1-${branch}.svc.cluster.local |                                      |
| Internal port    | 8080/tcp                                |                                      |
| External address | ---                                     |                                      |
| External port    | ---                                     |                                      |
| HTTP route       | /api                                     |                                      |