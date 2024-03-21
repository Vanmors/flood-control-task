# Ход мыслей

Когда нужно было выбрать хранилище данных,
были мысли на счёт хранения времени и количества запросов в
памяти приложения, но так как флуд-контроль может быть запущен
на нескольких экземплярах приложения одновременно нам нужна общая
база данных. Дальше мысли пали на реляционную базу данных, которая
обеспечит надёжное хранение наших данных, но после я понял
что в данной ситуации скорость важней чем надёжность, поэтому
мой выбор пал на Redis.

Изначально была идея реализовать это с помощью хранения map в Redis,
где в качестве ключа был бы идентификатор пользователя, а в значении хранилось бы время
послденего запроса и количество запросов за последний промежуток времени.
Чуть позже стало понятно, что это решение не со всем верно, так как мы не считаем
отклонённый запрос за вызов метода Check, а по условию это нужно учитывать.

Для реализации данной логики в флуд-контроле с
использованием Redis мы используем сортированный набор
(sorted set), в котором ключами будут идентификаторы пользователей,
а значениями будет время последнего вызова функции. Затем, для
определения количества вызовов за последние N секунд, мы
используем команду ZREMRANGEBYSCORE для удаления всех элементов,
время которых находится за пределами интервала последних
N секунд, после получаем все временные метки запросов за последние N секунд
и проверяем меньше ли их чем K. Дальше в любом случае сохраняем то, что метод был вызван
и возвращаем результат проверки.


Когда завершите задачу, в этом README опишите свой ход мыслей: как вы пришли к решению, какие были варианты и почему выбрали именно этот.

# Что нужно сделать

Реализовать интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

- Интерфейс FloodControl располагается в файле main.go.

- Флуд-контроль может быть запущен на нескольких экземплярах приложения одновременно, поэтому нужно предусмотреть общее хранилище данных. Допустимо использовать любое на ваше усмотрение.

# Необязательно, но было бы круто

Хорошо, если добавите поддержку конфигурации итоговой реализации. Параметры — на ваше усмотрение.
