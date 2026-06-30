begin;
insert into users (id, email)
values ('019f12cb-b476-77d1-9317-24b1b19aa6ca', 'karol@moroz.dev')
on conflict (email) do nothing;

insert into tunnels (id, subdomain, username, password_hash, user_id)
-- 35e0d095:f01494ad
-- a013612c:8b4ec2b5
-- fabf85c2:8abd93e2
values
    ('019f0610-fba4-708e-a9e7-40382b4f9f72', 'atrocious-jaguar', '35e0d095',
        '$argon2id$v=19$m=65536,t=3,p=4$6Cig8i4z9lr0sVhNra0qPg$sq4U2KMep1fPBh4YdSebUHYLxNoL642ff8ivPrzR6gU',
        '019f12cb-b476-77d1-9317-24b1b19aa6ca'),
    ('019f17f5-b463-7633-a829-d05990df55c8', 'voluptuous-primate', 'a013612c',
        '$argon2id$v=19$m=65536,t=3,p=4$VuLpxiW+HJCD1qibaRVJ+Q$eAF5IyCaTOsLAYtzTCUXhNRhHgBT8td83FM2GoMQP6A',
         '019f12cb-b476-77d1-9317-24b1b19aa6ca'),
    ('019f188b-fb6d-70b9-8f76-caf49726d748', 'succulent-chinese-meal', 'fabf85c2',
        '$argon2id$v=19$m=65536,t=3,p=4$OqINTVoaC901cdyzNzms+w$GD2Wp51hK0I1IPnXhEL/hSEBPnVDXZ4qwDGlZROT0aI',
        '019f12cb-b476-77d1-9317-24b1b19aa6ca')
on conflict (username) do nothing;

commit;
