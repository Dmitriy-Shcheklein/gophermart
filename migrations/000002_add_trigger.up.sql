create or replace function after_order_update()
    returns trigger as $$
        begin
            update balances
            set current = current + new.accrual - coalesce(old.accrual, 0)
            where user_id = new.user_id;
        return new;
    end;
$$ language plpgsql;

create trigger trg_after_order_update
after update on orders
for each row
execute function after_order_update();