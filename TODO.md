## The Grand List of TODOs

* SQL-based functions that modify data need to run an additional read query. For example, 
EntryService.CreateEntry inserts the row and then reads back the same row in the return statement. This is because the Ent insert call
does not return an updated version of the row. See https://github.com/ent/ent/issues/1131
* SQL-based services need unit testing