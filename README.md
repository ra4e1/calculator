# Инструкция

## Описание
Это калькулятор, который был сделан как контрольная задача второго спринта Яндекс Лицея на курсе по изучению GO. 

## Установка, запуск, остановка
Для **установки** калькулятора перейдите в папку, в которую необходимо установить приложение, и запустите следующие команды: 
```bash
git clone https://github.com/ra4e1/calculator.git

cd calculator

go mod vendor
```

Для **запуска** калькулятора в командную строку вводится:
```bash
go run internal/main.go
```
Для **остановки** программы в командную строку вводится сочетание клавиш "Ctrl"+"C".

## Пример использования

В качестве примера рассчитаем выражение: 2+2х2

1. Запустите программу.
1. В браузере откройте ссылку "http://localhost:8081/calc?q=2%2b2*2". В ответ будет выведен идентификатор запроса. Запомните его.
1. Перейдите по ссылке "http://localhost:8081/list" Вы должны увидеть ваш запрос в списке выполнения.
1. Для получения результата перейдите по ссылке "http://localhost:8081/answer?id=1" (используйте другой id если необходимо). Вы должны увидеть результат расчета или текст "ожидайте решения", если ответ еще не посчитан.
1. Остановите приложение и запустите его заново.
1. Откройте ссылку "http://localhost:8081/list". Вы должны увидеть тот же список запросов, что и до остановки приложения. Состояние и список запросов восстановились после перезапуска.

## Описание API

Запросы отправляются методом GET. Результат выдается в теле ответа.

### /calc
Запускает решение примера и возвращает идентификатор.

**Параметры**:
  * **q** (тип: string) - выражение для расчета, обязательный параметр. Должно содержать только цифры и знаки: ```+ - * / ( )```  

**Ответ**:  
* **requestID** - идентификатор запроса. Используется в дальнейшем в ```/answer``` для получения результата

Параметры передаются в URL и должны быть закодированы в соответствии с правилами кодировки url.

### /answer
Используется для получения результата расчета выражения, отправленно с помощью ```/calc```

**Параметры**:
  * **id** (тип: int) - идентификатор запроса, обязательный параметр  

**Ответ**:  
* **result** - результат  расчета выражения или следующие ошибки:

1. **ожидайте решения** - на данный момент выражение решается.
1. **отменено** - при решении выражения сервер был остановлен. Расчет не был завершен.
1. **Unexpected end of expression** - в выражении содержится ошибка.
1. **Cannot transition token types from NUMERIC [ ] to NUMERIC [ ]** - в выражении используются символы, которые программа не поддерживает.
1. **+Inf** - в выражении имеется деление на ноль.


### /list
Используется для просмотра списка всех выражений. Показывает: 
* идентификатор запроса
* выражение для расчета
* статус задачи (решено, решается, ошибка)
* ошибку, которая возникла при подсчетах.

Входные параметры отсутствуют.

## Конфигурация приложения

В приложении предусмотрены следующие конфигурационные параметры:
 * CALC_APP_PORT - порт для запуска http-сервера. По умолчанию: 8081
 * CALC_APP_DELAY - время задержки расчета выражения в секундах. Используется для замедления расчета, чтобы имитировать сложное вычисление. По умолчанию: 10 секунд.

Параметры приложения задаются через переменные окружения.

## Мониторинг воркеров
Расчет каждого выражения происходит в горутине. Список выполняющихся горутин можно увидеть с помощью функции ```/list```.

## Описание структуры приложения

У приложения есть следующие логические компоненты:

* Структура ```Application```. Конфигурирует и запускает приложение. Содержит одну функцию ```Run```
* Сервисы. Предназначены для выполнения какой-то определенной задачи. Есть два сервиса:
    * ```CalculatorService``` - рассчитывает выражения
    * ```StateService``` - сохраняет и восстанавливает состояние приложения 
* Структура ```Webserver```. Отвечает за запуск и работу http-сервера. Содержит обработчики http-запросов.