-- TODO нужно удалить файл, если у вас нет функций
create or replace function foo()
    returns void language plpgsql as $$
begin
    raise notice 'foo()';
end;
$$;
