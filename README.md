# converter
converter это простое консольное приложение для преобразования файлов из одного формата в другой. Сделана исключительно в тестовых целях.

# Overview
Поддерживаемые команды:
```
conveter conver --source=source --destination=destination [--prettyprint] [overwrite]
conveter version
conveter help
```
Примеры:
```
./converter convert --source=/home/www/Work/my/converter/test/docker-compose.yml --destination=/home/www/Work/my/converter/test2/docker-compose.json --prettyprint --overwrite
./converter convert --source=/home/www/Work/my/converter/test --destination=/home/www/Work/my/converter/test2 --prettyprint
```

# Install
```
make build GOROOT=goorot
```
Пример:
```
make build GOROOT='/home/www/apps/go/go1.16.5'
```
