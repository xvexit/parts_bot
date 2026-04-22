-- Изменено: добавлено "UNIQUE (name, brand)"
CREATE TABLE IF NOT EXISTS catalog_parts (
    part_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(120) NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    delivery_day INT NOT NULL CHECK (delivery_day >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (name, brand) -- <-- ВОТ ИСПРАВЛЕНИЕ
);

-- 1. CTE с базовыми ценами и брендами
WITH parts_base AS (
    SELECT
        t.name,
        unnest(t.brands) AS brand, -- "Разворачиваем" массив брендов
        t.base_price -- Используем базовую цену как основу
    FROM (VALUES
        -- Формат: (Название, {Бренды}, Базовая_цена_руб)
        -- ФИЛЬТРЫ
        ('Фильтр масляный', ARRAY['MANN','Filtron','Mahle'], 600),
        ('Фильтр воздушный', ARRAY['MANN','Filtron','Mahle'], 900),
        ('Фильтр салонный', ARRAY['MANN','Filtron','Mahle'], 1100),
        ('Фильтр топливный', ARRAY['MANN','Bosch','Mahle'], 2200),
        -- ЗАЖИГАНИЕ / ТОПЛИВО
        ('Свеча зажигания', ARRAY['NGK','Bosch'], 700),
        ('Свеча накаливания', ARRAY['Bosch','NGK'], 1600),
        ('Катушка зажигания', ARRAY['Bosch'], 4500),
        ('Форсунка топливная', ARRAY['Bosch','Denso'], 11000),
        ('Топливный насос', ARRAY['Bosch','Pierburg'], 9500),
        ('Регулятор давления топлива', ARRAY['Bosch'], 3500),
        -- ТОРМОЗА
        ('Колодки тормозные передние', ARRAY['ATE','TRW','Brembo'], 4800),
        ('Колодки тормозные задние', ARRAY['ATE','TRW','Brembo'], 3500),
        ('Диск тормозной', ARRAY['ATE','Brembo'], 6500),
        ('Барабан тормозной', ARRAY['ATE','TRW'], 5500),
        ('Суппорт тормозной', ARRAY['ATE','Brembo'], 13000),
        ('Ремкомплект суппорта', ARRAY['ATE','TRW'], 1500),
        ('Тормозной цилиндр', ARRAY['ATE','TRW'], 2500),
        ('Главный тормозной цилиндр', ARRAY['ATE','TRW'], 8000),
        ('Вакуумный усилитель тормозов', ARRAY['ATE'], 11000),
        ('Тормозной шланг', ARRAY['ATE','TRW'], 1800),
        -- ПОДВЕСКА
        ('Амортизатор передний', ARRAY['Sachs','KYB'], 8000),
        ('Амортизатор задний', ARRAY['Sachs','KYB'], 7200),
        ('Пружина подвески', ARRAY['Lesjofors','KYB'], 4000),
        ('Стойка стабилизатора', ARRAY['TRW','Lemforder'], 1900),
        ('Втулка стабилизатора', ARRAY['Febi','Lemforder'], 600),
        ('Рычаг подвески', ARRAY['TRW','Lemforder'], 9500),
        ('Шаровая опора', ARRAY['TRW','Lemforder'], 2400),
        ('Сайлентблок', ARRAY['Febi','Lemforder'], 900),
        ('Опора двигателя', ARRAY['Febi','Lemforder'], 4500),
        ('Опора коробки передач', ARRAY['Febi','Lemforder'], 4200),
        -- РЕМНИ / ГРМ
        ('Ремень генератора', ARRAY['Contitech','Gates'], 900),
        ('Ремень ГРМ', ARRAY['Contitech','Gates'], 2500),
        ('Ремень ГРМ комплект', ARRAY['Contitech','Gates'], 11000),
        ('Ролик натяжной', ARRAY['INA','SKF'], 3500),
        ('Ролик обводной', ARRAY['INA','SKF'], 3200),
        ('Шкив коленвала', ARRAY['INA','Febi'], 7000),
        ('Шкив генератора', ARRAY['INA','Febi'], 5500),
        ('Цепь ГРМ', ARRAY['INA'], 9000),
        ('Натяжитель цепи', ARRAY['INA'], 6500),
        ('Успокоитель цепи', ARRAY['INA'], 4000),
        -- ОХЛАЖДЕНИЕ
        ('Помпа охлаждения', ARRAY['SKF','HEPU'], 5500),
        ('Радиатор охлаждения', ARRAY['Nissens','Valeo'], 12000),
        ('Радиатор кондиционера', ARRAY['Nissens','Valeo'], 11000),
        ('Интеркулер', ARRAY['Nissens'], 13000),
        ('Термостат', ARRAY['Mahle','Behr'], 2800),
        ('Датчик температуры', ARRAY['Bosch'], 1200),
        ('Вентилятор охлаждения', ARRAY['Valeo'], 9000),
        ('Крышка радиатора', ARRAY['Febi'], 500),
        ('Расширительный бачок', ARRAY['Febi'], 2500),
        ('Патрубок охлаждения', ARRAY['Gates'], 1500),
        -- ЭЛЕКТРИКА
        ('Аккумулятор', ARRAY['Varta','Bosch'], 9500),
        ('Генератор', ARRAY['Bosch','Valeo'], 28000),
        ('Стартер', ARRAY['Bosch','Valeo'], 19000),
        ('Реле стартера', ARRAY['Bosch'], 800),
        ('Предохранитель', ARRAY['Bosch'], 100),
        ('Блок предохранителей', ARRAY['Bosch'], 7000),
        ('Датчик ABS', ARRAY['Bosch'], 3500),
        ('Датчик коленвала', ARRAY['Bosch'], 2800),
        ('Датчик распредвала', ARRAY['Bosch'], 2700),
        ('Датчик кислорода', ARRAY['Bosch'], 6500),
        -- СВЕТ / КУЗОВ
        ('Фара передняя', ARRAY['Valeo','Hella'], 25000),
        ('Фонарь задний', ARRAY['Valeo','Hella'], 12000),
        ('Противотуманная фара', ARRAY['Valeo','Hella'], 7000),
        ('Лампа галогенная', ARRAY['Osram','Philips'], 300),
        ('Лампа LED', ARRAY['Osram','Philips'], 1500),
        ('Стекло фары', ARRAY['Depo'], 3000),
        ('Мотор стеклоочистителя', ARRAY['Bosch'], 5500),
        ('Щетка стеклоочистителя', ARRAY['Bosch','Valeo'], 1800),
        ('Бачок омывателя', ARRAY['Febi'], 2000),
        ('Насос омывателя', ARRAY['Bosch'], 1300),
        -- ВЫХЛОП
        ('Глушитель', ARRAY['Bosal'], 9000),
        ('Резонатор', ARRAY['Bosal'], 6000),
        ('Катализатор', ARRAY['Bosal'], 35000),
        ('Лямбда-зонд', ARRAY['Bosch'], 6500),
        ('Прокладка выпускного коллектора', ARRAY['Elring'], 700),
        ('Прокладка ГБЦ', ARRAY['Elring'], 2500),
        ('Прокладка клапанной крышки', ARRAY['Elring'], 1200),
        ('Сальник коленвала', ARRAY['Elring'], 900),
        ('Сальник распредвала', ARRAY['Elring'], 800),
        ('Комплект прокладок двигателя', ARRAY['Elring'], 8000),
        -- ТРАНСМИССИЯ
        ('Коробка передач', ARRAY['ZF'], 150000),
        ('Сцепление комплект', ARRAY['LUK','Sachs'], 17000),
        ('Диск сцепления', ARRAY['LUK','Sachs'], 8000),
        ('Корзина сцепления', ARRAY['LUK','Sachs'], 9000),
        ('Выжимной подшипник', ARRAY['SKF','INA'], 4500),
        ('Маховик', ARRAY['LUK'], 22000),
        ('Приводной вал', ARRAY['GKN'], 18000),
        ('ШРУС', ARRAY['GKN'], 8500),
        ('Пыльник ШРУСа', ARRAY['Febi'], 1200),
        ('Редуктор', ARRAY['ZF'], 80000),
        -- РУЛЕВОЕ
        ('Рулевая рейка', ARRAY['ZF','TRW'], 45000),
        ('Наконечник рулевой тяги', ARRAY['TRW','Lemforder'], 2200),
        ('Рулевая тяга', ARRAY['TRW','Lemforder'], 3200),
        ('Гидроусилитель руля', ARRAY['ZF'], 25000),
        ('Насос ГУР', ARRAY['ZF'], 18000),
        ('Жидкость ГУР', ARRAY['Febi'], 1000),
        ('Подушка безопасности', ARRAY['Bosch'], 25000),
        ('Замок зажигания', ARRAY['Bosch'], 6000),
        ('Кнопка старт-стоп', ARRAY['Bosch'], 4000),
        ('Иммобилайзер', ARRAY['Bosch'], 12000)
    ) AS t(name, brands, base_price)
),

-- 2. CTE для расчета финальной цены с учетом бренда
parts_final_price AS (
    SELECT
        name,
        brand,
        -- Применяем коэффициент в зависимости от бренда
        (base_price *
            CASE
                WHEN brand IN ('Brembo', 'ATE', 'Lemforder', 'Sachs', 'LUK', 'GKN', 'ZF', 'Hella', 'Valeo', 'Behr') THEN 1.25 -- Премиум +25%
                WHEN brand IN ('Bosch', 'NGK', 'SKF', 'INA', 'Contitech', 'Gates', 'Mahle', 'Nissens', 'Pierburg') THEN 1.1 -- Средний +10%
                ELSE 0.95 -- Бюджетный -5%
            END
        )::BIGINT as calculated_price
    FROM parts_base
)

-- 3. Вставляем данные в таблицу
INSERT INTO catalog_parts (part_id, name, brand, price, delivery_day)
SELECT
    'A-' || lpad(row_number() OVER (ORDER BY p.name, p.brand)::text, 4, '0'),
    p.name,
    p.brand,
    -- Добавляем небольшую случайную "погрешность" (+-5%) к расчетной цене для живости
    (p.calculated_price * (0.95 + random() * 0.1))::BIGINT,
    (1 + floor(random() * 5))::INT
FROM parts_final_price p
-- Гарантируем, что при повторном запуске дубликаты не создадутся
ON CONFLICT (name, brand) DO NOTHING;
