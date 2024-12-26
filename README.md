Приложение ISBN Locator позволяет:
  - получать информацию о книге по ISBN,
  - добавлять информацию в хранилище,
  - обновлять и удалять информацию.

Запуск и остановка приложения:

```shell
make start
```
```shell
make stop
```

Использование:

```
curl -i -X POST localhost:8000/books --data '{"ISBN":"2-266-11156-6","Title":"Mybook1","Author":"Author1","Year":"2021"}'

curl -i -X PUT localhost:8000/books/2-266-11156-6 --data '{"ISBN":"2-266-11156-6","Title":"Mybook2","Author":"Author2","Year":"2022"}'

curl -i -X GET localhost:8000/books/2-266-11156-6

curl -i -X DELETE localhost:8000/books/2-266-11156-6
```

Принимает международный книжный номер в стандартах ISBN-10 и ISBN-13. Например: `2-266-11156-6, 0-712-67263-X, 978-5-907488-10-6`.

Notes:

- Standard Go Layout: https://github.com/golang-standards/project-layout/
- Another example: https://github.com/herryg91/go-clean-architecture/tree/main/examples/book-rest-api
- https://github.com/herryg91/go-clean-architecture/blob/main/README.md#folder-structure
