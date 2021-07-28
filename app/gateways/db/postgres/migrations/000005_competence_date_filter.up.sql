begin;

create index if not exists idx_entry_competence_date
    on entry using btree (competence_date);

commit;
