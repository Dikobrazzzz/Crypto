# Отслеживание балансов и транзакций криптокошельков

0. Поставить golang
1. Делаем api: 
    1. Написать hello world
    2. По ручкам: (делаем пока с заглушками и без похода в базу)
        Используй gin framework
        1. `POST` `/address` request body{"wallet_address", "chain_name", "crypto_name", "tag", ...}, response 200 body {"id", "wallet_address", "chain_name",  "crypto_name", "balance", "tag", ...}, если была internal error -> 500, error
        2. `GET` `/address/:id` response 200 body {"id", "wallet_address", "chain_name", "tag", ...}, если была internal error -> 500, error
        3. ручка для получения всех кошельков
        4. `PUT` `/address/tag` request body{"id", "tag"}, response: 200, "id"

2. слайс vs массив, стурктура слайса, что происходит при append
3. map, бакеты, миграции, коллизии