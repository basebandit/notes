# Notes

Notes is a simple notes api.

## Getting Started

- Create a new empty database (MySQL оr PostgreSQL) that will be used as a data store and a database user with all privileges granted on this database.

* Download and compile the notes binary:

  ```
    $ go get -u github.com/basebandit/notes
  ```

* Inside an empty directory run:

  ```
  $ ./notes init
  ```

  This will generate an initial configuration file "notes.conf" inside the current dir.
  Edit the configuration file to set the server listen address, the base url, the database, etc.

- Run the following command to start the notes api server.
  ```
  $ ./notes start
  ```
