CREATE DATABASE source;

CREATE USER source;

ALTER USER source WITH PASSWORD 'source';

\c source

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

GRANT SELECT ON users TO source;

INSERT INTO
    users (name, email)
VALUES
    ('Alan Gou', 'alan.gou@source.co'),
    ('Bob Francis', 'bob.francis@source.co'),
    ('Cathy Lee', 'cathy.lee@source.co'),
    ('David Chen', 'david.chen@source.co'),
    ('Alice Zhang', 'alice.zhang@source.co'),
    ('Brian Lee', 'brian.lee@source.co'),
    ('Cindy Wang', 'cindy.wang@source.co'),
    ('Daniel Kim', 'daniel.kim@source.co'),
    ('Amy Chen', 'amy.chen@source.co'),
    ('Brandon Lee', 'brandon.lee@source.co'),
    ('Christine Kim', 'christine.kim@source.co'),
    ('Derek Chen', 'derek.chen@source.co'),
    ('Anna Zhang', 'anna.zhang@source.co'),
    ('Benjamin Lee', 'benjamin.lee@source.co'),
    ('Claire Wang', 'claire.wang@source.co'),
    ('Edward Kim', 'edward.kim@source.co'),
    ('Amanda Chen', 'amanda.chen@source.co'),
    ('Frank Lee', 'frank.lee@source.co'),
    ('Grace Kim', 'grace.kim@source.co'),
    ('Ethan Chen', 'ethan.chen@source.co'),
    ('Hannah Zhang', 'hannah.zhang@source.co'),
    ('Isaac Lee', 'isaac.lee@source.co'),
    ('Jessica Wang', 'jessica.wang@source.co'),
    ('George Kim', 'george.kim@source.co'),
    ('Karen Chen', 'karen.chen@source.co'),
    ('Henry Lee', 'henry.lee@source.co'),
    ('Linda Zhang', 'linda.zhang@source.co'),
    ('Jack Chen', 'jack.chen@source.co'),
    ('Mary Lee', 'mary.lee@source.co'),
    ('Nancy Wang', 'nancy.wang@source.co'),
    ('Adam Kowalski', 'adam.kowalski@source.co'),
    ('Barbara Nowak', 'barbara.nowak@source.co'),
    ('Cezary Wojcik', 'cezary.wojcik@source.co'),
    ('Dorota Kaczmarek', 'dorota.kaczmarek@source.co'),
    (
        'Andrzej Zielinski',
        'andrzej.zielinski@source.co'
    ),
    ('Beata Szymanska', 'beata.szymanska@source.co'),
    ('Czeslaw Wozniak', 'czeslaw.wozniak@source.co'),
    (
        'Dominika Kowalczyk',
        'dominika.kowalczyk@source.co'
    ),
    (
        'Artur Lewandowski',
        'artur.lewandowski@source.co'
    ),
    ('Bozena Jankowska', 'bozena.jankowska@source.co'),
    ('Dariusz Mazur', 'dariusz.mazur@source.co'),
    (
        'Alicja Wojciechowska',
        'alicja.wojciechowska@source.co'
    ),
    ('Bartosz Krawczyk', 'bartosz.krawczyk@source.co'),
    ('Celina Pawlowska', 'celina.pawlowska@source.co'),
    ('Damian Michalski', 'damian.michalski@source.co'),
    ('Elzbieta Zajac', 'elzbieta.zajac@source.co'),
    ('Filip Gorski', 'filip.gorski@source.co'),
    (
        'Gabriela Tomaszewska',
        'gabriela.tomaszewska@source.co'
    ),
    (
        'Henryk Wasilewski',
        'henryk.wasilewski@source.co'
    ),
    (
        'Izabela Kaczmarczyk',
        'izabela.kaczmarczyk@source.co'
    ),
    ('Janusz Krol', 'janusz.krol@source.co'),
    (
        'Katarzyna Wieczorek',
        'katarzyna.wieczorek@source.co'
    ),
    ('Lukasz Nowicki', 'lukasz.nowicki@source.co'),
    (
        'Magdalena Kowalewska',
        'magdalena.kowalewska@source.co'
    ),
    (
        'Norbert Zielinski',
        'norbert.zielinski@source.co'
    ),
    ('Olga Sadowska', 'olga.sadowska@source.co'),
    ('Piotr Kaczynski', 'piotr.kaczynski@source.co'),
    ('Renata Kowalczyk', 'renata.kowalczyk@source.co'),
    (
        'Sebastian Wojciechowski',
        'sebastian.wojciechowski@source.co'
    ),
    (
        'Teresa Nowakowska',
        'teresa.nowakowska@source.co'
    ),
    ('Urszula Kowalska', 'urszula.kowalska@source.co'),
    (
        'Wojciech Zielinski',
        'wojciech.zielinski@source.co'
    ),
    ('Zofia Krol', 'zofia.krol@source.co'),
    ('Antonio Rossi', 'antonio.rossi@source.co'),
    ('Giuseppe Russo', 'giuseppe.russo@source.co'),
    ('Marco Ferrari', 'marco.ferrari@source.co'),
    ('Paolo Bianchi', 'paolo.bianchi@source.co'),
    ('Simone Romano', 'simone.romano@source.co'),
    ('Avery Johnson', 'avery.johnson@source.co'),
    ('Blake Thompson', 'blake.thompson@source.co'),
    ('Cameron Davis', 'cameron.davis@source.co'),
    ('Dylan Wilson', 'dylan.wilson@source.co'),
    ('Ethan Brown', 'ethan.brown@source.co'),
    ('Faith Taylor', 'faith.taylor@source.co'),
    ('Gavin Anderson', 'gavin.anderson@source.co'),
    ('Haley Martin', 'haley.martin@source.co'),
    ('Isabella White', 'isabella.white@source.co'),
    ('Jacob Clark', 'jacob.clark@source.co'),
    ('Kaitlyn Wright', 'kaitlyn.wright@source.co'),
    ('Landon Scott', 'landon.scott@source.co'),
    ('Mia Green', 'mia.green@source.co'),
    ('Nathan Baker', 'nathan.baker@source.co'),
    ('Olivia King', 'olivia.king@source.co'),
    ('Parker Adams', 'parker.adams@source.co'),
    ('Quinn Evans', 'quinn.evans@source.co'),
    ('Riley Cooper', 'riley.cooper@source.co'),
    ('Samantha Parker', 'samantha.parker@source.co'),
    ('Tyler Turner', 'tyler.turner@source.co'),
    ('Victoria Collins', 'victoria.collins@source.co'),
    ('Wyatt Murphy', 'wyatt.murphy@source.co'),
    ('Xavier Mitchell', 'xavier.mitchell@source.co'),
    ('Yasmine Perez', 'yasmine.perez@source.co'),
    ('Zachary Nelson', 'zachary.nelson@source.co');