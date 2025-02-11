<p align="center"><a href="https://www.codechefvit.com" target="_blank"><img src="https://i.ibb.co/4J9LXxS/cclogo.png" width=160 title="CodeChef-VIT" alt="Codechef-VIT"></a>
</p>

# Devsoc Backend '25

The official Backend API for DEVSOC'25 Hackathon Portal


## Features
- Fully-featured REST API built with Go
- Postgres integration using SQLC for typesafe bindings
- Structured routing and controllers 
- JWT authentication and middleware
- Comprehensive logging and error handling
- Docker and docker-compose configuration for containerized deployment
- API documentation generation via scalar



## Tech Stack
 - [Go](https://golang.org/)
    - [Echo](https://echo.labstack.com/): Minimalist Go web framework
    - [SQLC](https://sqlc.dev/): Type-safe Go from SQL
    - [pgx](https://github.com/jackc/pgx): PostgreSQL driver for Go
    - [gomail](https://github.com/go-gomail/gomail): Send emails in Go
    - [go-redis](https://github.com/go-redis/redis): Redis client for Go
 - [PostgreSQL](https://www.postgresql.org/): Open-source relational database
 - [Docker](https://www.docker.com/): Container platform
 - [Redis](https://redis.io/): In-memory data store
 - [Goose](https://github.com/pressly/goose): DB migration tool for Go

## How To Run

#### Clone the repo
```sh
$ git clone https://github.com/CodeChefVIT/devsoc-be-25
$ cd devsoc-be-25
```

#### Download Dependencies
```sh
$ go mod download
```

#### Fill Environment Variables
```sh
$ cp .env.sample .env
```

#### Start Containers
```sh
$ docker compose up --build
```

#### Apply Migrations
```sh
$ make up
```
## Contributors
<table>
	<tr align="center" style="font-weight:bold">
<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/71623796?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Vedant Matanhelia</p>
	<p align="center">
		<a href = "https://github.com/RustyDev24"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/vedant-matanhelia-aa171027b/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>

<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/155614230?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Soham Mahapatra</p>
	<p align="center">
		<a href = "https://github.com/Soham-Maha"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/soham-mahapatra-433bb428a/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>

<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/86644389?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Aman L</p>
	<p align="center">
		<a href = "https://github.com/Killerrekt"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/aman-l-922819251/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>

<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/139199971?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Vaibhav Sijaria</p>
	<p align="center">
		<a href = "https://github.com/sophic00"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/vaibhav-sijaria/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>






</tr>

<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/74227363?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Jothish Kamal</p>
	<p align="center">
		<a href = "https://github.com/JothishKamal"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/jothishkamal/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>

<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/80804989?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Aman Singh</p>
	<p align="center">
		<a href = "https://github.com/DevloperAmanSingh"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/amansingh2112/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>

<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/140488187?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Abhinav Anand</p>
	<p align="center">
		<a href = "https://github.com/Abhinav-055"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/abhinav-anand--/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>

<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/56132559?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Abhinav Ganeshan</p>
	<p align="center">
		<a href = "https://github.com/Abh1noob"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/abhinav-gk/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>
</tr>
<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/143117260?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Abhinav Garg</p>
	<p align="center">
		<a href = "https://github.com/ABHINAVGARG05"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/abhinav-garg-75798028a/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>

<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/67090539?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Abhinav Pant</p>
	<p align="center">
		<a href = "https://github.com/abhitrueprogrammer"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/abhinav-pant-081b79243/">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>
<td>
	<p align="center">
		<img src = "https://avatars.githubusercontent.com/u/84934511?v=4" width="200" height="200" alt="profilepic" style="border: 2px solid grey; width: 170px; height:170px">
	</p>
	<p style="font-size:17px; font-weight:600;">Nishant Gupta</p>
	<p align="center">
		<a href = "https://github.com/NishantGupt786"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
		<a href = "https://www.linkedin.com/in/nishant-gupta-12913221b">
			<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
		</a>
	</p>
</td>
	</tr>
</table>

# License

Copyright © 2025, [CodeChef-VIT](https://github.com/CodeChefVIT) and all other contributors.
Released under the [MIT License](LICENSE).

<p align="center">
Made with :heart: by <a href="https://www.codechefvit.com" target="_blank">CodeChef-VIT</a>
</p>