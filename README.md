# go_qm - QueryMongo

Данный пакет преобразует запрос Query в запрос в MongoDB. Вы можете фильтровать, сортировать, указывать лимиты и страницы. Для того чтобы все работало нужно передать на вход просто Query параметры в формате `query map[string]string`.

* `query map[string]string` - это карта по ключу и значению фильтрации/вывода

## Установка

```go
go get github.com/shinpi-tech/go_qm
```

## Как пользоваться

Для выполнения фильтрации определенных типов данных таких как: число и дата. Используются оператор `~`.

Пример:
```
?price=5000~6000&available=true&params.load_index=12-12-2024~!29-12-2024&params.width=175,195,205&params.width=175&brand=6302ac8a85bafafe377bd7dd,6302ac8a85bafafe377bd7dd&sort=-price,name&page=2&limit=30
```



