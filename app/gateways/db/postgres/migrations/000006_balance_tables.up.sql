begin;

drop table if exists account_balance;
drop table if exists aggregated_query_balance;

drop function if exists _get_account_balance;
drop function if exists _get_account_balance_since;
drop procedure if exists _update_account_balance;
drop procedure if exists _insert_account_balance;
drop function if exists get_account_balance;

drop function if exists _get_aggregated_query_balance;
drop function if exists _get_aggregated_query_balance_since;
drop procedure if exists _update_aggregated_query_balance;
drop procedure if exists _insert_aggregated_query_balance;
drop function if exists query_aggregated_account_balance;

create table if not exists account_balance
(
    balance bigint      not null,
    tx_date timestamptz not null,
    account text        primary key
) with (fillfactor = 70);

create or replace procedure _update_account_balance(
    _account text, _balance bigint, _dt timestamptz
)
    language sql
as
$$
    update account_balance
    set
        balance = _balance,
        tx_date = _dt
    where account = _account;
$$;

create or replace procedure _insert_account_balance(
    _account text, _balance bigint, _dt timestamptz
)
    language sql
as
$$
    insert into account_balance (balance, tx_date, account)
    values (_balance, _dt, _account);
$$;

--
-- Analytic account
--

create or replace function _get_analytic_account_balance(_account ltree)
    returns table
        (
            partial_balance bigint,
            partial_date    timestamptz,
            recent_balance  bigint,
            recent_version  int
        )
    language sql
as
$$
    select
        sum(sub.balance) filter (where sub.row_number > 1) as partial_balance,
        max(created_at)  filter (where sub.row_number = 2) as partial_date,

        sum(sub.balance) filter (where sub.row_number = 1) as recent_balance,
        max(version)     filter (where sub.row_number = 1) as recent_version
    from (
        select
            coalesce(sum(amount) filter (where operation = 1), 0) -
            coalesce(sum(amount) filter (where operation = 2), 0) as balance,
            max(version) as version,
            created_at,
            row_number() over (order by created_at desc) as row_number
        from entry
        where account = _account
        group by created_at
        order by created_at desc
    ) sub
$$ stable rows 1;

create or replace function _get_analytic_account_balance_since(_account ltree, _dt timestamptz)
    returns table
        (
            partial_balance bigint,
            partial_date    timestamptz,
            recent_balance  bigint,
            recent_version  int
        )
    language sql
as
$$
    select
        sum(sub.balance) filter (where sub.row_number > 1) as partial_balance,
        max(created_at)  filter (where sub.row_number = 2) as partial_date,

        sum(sub.balance) filter (where sub.row_number = 1) as recent_balance,
        max(version)     filter (where sub.row_number = 1) as recent_version
    from (
        select
            coalesce(sum(amount) filter (where operation = 1), 0) -
            coalesce(sum(amount) filter (where operation = 2), 0) as balance,
            max(version) as version,
            created_at,
            row_number() over (order by created_at desc) as row_number
        from entry
        where
            account = _account
            and created_at > _dt
        group by created_at
        order by created_at desc
    ) sub
$$ stable rows 1;

create or replace function get_analytic_account_balance(
    in _account ltree,
    out total_balance bigint, out version int
)
    returns record
    language plpgsql
as
$$
declare
    _existing_balance   bigint;
    _existing_date      timestamptz;

    _partial_balance    bigint;
    _partial_date       timestamptz;
begin
    select
        balance,
        tx_date
    into
        _existing_balance,
        _existing_date
    from
        account_balance
    where
        account = _account::text;

    if (_existing_balance is null) then
        select
            partial_balance,
            partial_date,
            coalesce(partial_balance, 0) + recent_balance,
            recent_version
        into
            _partial_balance,
            _partial_date,
            total_balance,
            version
        from
            _get_analytic_account_balance(_account);

        -- No entries found for the given account
        if (version is null) then
            raise no_data_found;
        -- Only recent balance exists, so return it without creating snapshot
        elsif (_partial_balance is null) then
            return;
        end if;

        call _insert_account_balance(
            _account => _account::text,
            _balance => _partial_balance,
            _dt => _partial_date
        );

        return;
    end if;

    select
        _existing_balance + partial_balance,
        partial_date,

        _existing_balance + coalesce(partial_balance, 0) + coalesce(recent_balance, 0),
        recent_version
    into
        _partial_balance,
        _partial_date,

        total_balance,
        version
    from
        _get_analytic_account_balance_since(_account, _existing_date);

    -- No new entries exists
    if (_partial_date is null) then
        return;
    end if;

    call _update_account_balance(
        _account => _account::text,
        _balance => _partial_balance,
        _dt => _partial_date
    );
end;
$$ volatile;

--
-- Synthetic account
--

create or replace function _get_synthetic_account_balance(_account lquery)
    returns table
        (
            partial_balance bigint,
            partial_date    timestamptz,
            recent_balance  bigint
        )
    language sql
as
$$
select sum(sub.balance) filter (where sub.row_number > 1) as partial_balance,
       max(created_at)  filter (where sub.row_number = 2) as partial_date,
       sum(sub.balance) filter (where sub.row_number = 1) as recent_balance
from (
         select
            coalesce(sum(amount) filter (where operation = 1), 0) -
            coalesce(sum(amount) filter (where operation = 2), 0) as balance,
            created_at,
            row_number() over (order by created_at desc) as row_number
         from entry
         where account ~ _account
         group by created_at
         order by created_at desc
     ) sub
$$ stable
   rows 1
;

create or replace function _get_synthetic_account_balance_since(_account lquery, _dt timestamptz)
    returns table
        (
            partial_balance bigint,
            partial_date    timestamptz,
            recent_balance  bigint
        )
    language sql
as
$$
select sum(sub.balance) filter (where sub.row_number > 1) as partial_balance,
       max(created_at) filter (where sub.row_number = 2)  as partial_date,
       sum(sub.balance) filter (where sub.row_number = 1) as recent_credit
from (
         select
            coalesce(sum(amount) filter (where operation = 1), 0) -
            coalesce(sum(amount) filter (where operation = 2), 0) as balance,
            created_at,
            row_number() over (order by created_at desc) as row_number
         from entry
         where account ~ _account
           and created_at > _dt
         group by created_at
         order by created_at desc
     ) sub
$$ stable
   rows 1
;

create or replace function get_synthetic_account_balance(
    in _account lquery, out total_balance bigint
)
    returns bigint
    language plpgsql
as
$$
declare
    _existing_balance bigint;
    _existing_date    timestamptz;
    _partial_balance  bigint;
    _partial_date     timestamptz;
begin
    select balance,
           tx_date
    into
        _existing_balance,
        _existing_date
    from account_balance
    where account = _account::text;

    if (_existing_balance is null) then
        select partial_balance,
               partial_date,
               coalesce(partial_balance, 0) + recent_balance
        into
            _partial_balance,
            _partial_date,
            total_balance
        from
            _get_synthetic_account_balance(_account);

        if (total_balance is null) then
            raise no_data_found;
        elsif (_partial_balance is null) then
            return;
        end if;

        call _insert_account_balance(
            _account => _account::text,
            _balance => _partial_balance,
            _dt => _partial_date
        );

        return;
    end if;

    select _existing_balance + partial_balance,
           partial_date,
           _existing_balance + coalesce(partial_balance, 0) + coalesce(recent_balance, 0)
    into
        _partial_balance,
        _partial_date,
        total_balance
    from
        _get_synthetic_account_balance_since(_account, _existing_date);

    if (_partial_date is null) then
        return;
    end if;

    call _update_account_balance(
        _account => _account::text,
        _balance => _partial_balance,
        _dt => _partial_date
    );
end;
$$ volatile;

commit;
