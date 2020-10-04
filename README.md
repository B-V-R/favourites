### System requirements:
    1. install latest version of docker
   
### How to run
    1. clone `https://github.com/B-V-R/favourites.git`
    2. change directory to project directory ex: `cd favourites`
    3. `docker-compose up --build --remove-orphans`
    4. open `http://localhost:80/`
    5. `Sign Up`
    6. `Sign In`
    7. Images will be displayed, click on `Add to Favourite`, `Add to Favourite` will turn into `light blue` color.
    8. Refresh browser tab and check.
    9. `Sgign out` (Below images) expires session, will remove images from favourites list
    
### Technologies used
    1. Backend code written in `Golang`.
    2. Frontend code written in `HTML`, `CSS` and `Javascript`.
    3. Database is `Postgres`

### Security
    1. User sign up and sign in are required to view images and add to favourites.
    2. User password is hashed with `bcrypt` (Faster than SHA and more secure)
   
    