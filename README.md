# What is the project?

## High level overview of how it works?

### Installation

1. git clone ...
2. open two terminals
3. in one terminal `cd face_recognition`
4. create a virtual python environment
5. `pip install -r requirements.txt`
6. in the other terminal `cd gnark`
7. create a .env inside gnark directory for your `DATABASE_URL=`,  `PROVING_KEY_PATH=`, `VERIFYING_KEY_PATH=`
8. create the postgres database
   1. `brew install psql-17`
   2. `CREATE DATABASE my_project_db;`
   3. `\c my_project_db`
   4. create this table
   5. put `postgres://yourUsername:yourPassword@localhost:yourPort/yourDB` as the `DATABASE_URL=`

   ```sql
   CREATE TABLE IF NOT EXISTS public.zetable (
    username    TEXT    NOT NULL,
    proof       BYTEA,
    helper_data JSONB,
    CONSTRAINT  zetable_username_key UNIQUE (username)
    ); 
    ```

#### Running the project

1. `cd gnark`
2. `go run .`
3. `cd face_recognition`
4. `python src/options.py enroll` for enrollment, then write a username
5. `python src/options.py verify` for verification, and write the same username

#### Contributors

I own both github accounts but misconfigured my global git config file.

#### Acknowledgements

- [gnark](https://github.com/ConsenSys/gnark)
- [deepface](https://github.com/serengil/deepface)
- [fuzzy_extractor](https://github.com/carter-yagemann/python-fuzzy-extractor)
