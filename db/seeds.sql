begin;

  insert into hosts (id, subdomain, username, password_hash)
  -- 35e0d095:f01494ad
  values ('019f0610-fba4-708e-a9e7-40382b4f9f72', 'atrocious-jaguar', '35e0d095', '$argon2id$v=19$m=65536,t=3,p=4$6Cig8i4z9lr0sVhNra0qPg$sq4U2KMep1fPBh4YdSebUHYLxNoL642ff8ivPrzR6gU')
  on conflict (username) do nothing;

commit;
