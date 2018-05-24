## Golang REST API starter

This go starter is here for you to take it and modify according to your needs.
This project promotes ideomatic go code.
And if there's any parts which has a bad design pattern,
expects you to provide feedback by opening a issue.


### What is this?
The idea is to have everything included (a.k.a. batteries included) starter project,
after which all you need to do would be to work on your idea,
instead of setting everything up.

### Backend
- MongoDB for data storage
- redis for caching
- docker for providing redis and MongoDB for development

### What features are included?
This project provides the following features:
- [x] JWT based Authentication
- [x] Pre configured endpoints (see handlers for more info)
- [x] Permissions
- [x] Middlewares (includes must_have_permission, rate_limit, auth_middleware etc)
- [x] Pre defined way of using Database with integration tests
- [x] Ability to cover 100% of code with tests
- [x] Ability to unit test each handlers
- [x] Transformers (so you can transform your db representation into response)
- [x] mocking (using mockgen)
- [x] migrations
- [x] redis cache
- [x] .env config
- [x] commands
- [x] travis and circleci support
- [x] docker-compose for database and redis
- [x] logging with logrus (and slack pre configured)
- [ ] easy email support
- [ ] sign up
- [ ] login with open auth provider
- [ ] recaptcha support (spam protection)
- [ ] open auth server
- [ ] command line tool to generate handlers, middlewares, etc
- [ ] limit login attempts
- [ ] side effects handling
- And more to come

### Packages Used
- [dep](https://github.com/golang/dep) for managing packages
- [gorilla mux](https://github.com/gorilla/mux) for routing
- [logrus](https://github.com/sirupsen/logrus) for logging
- [respond.v1](https://github.com/matryer/respond) for responding to requests
- [mgo](github.com/globalsign/mgo) for connecting to mongodb
- [go-redis](github.com/go-redis/redis) for redis
- [testify](github.com/stretchr/testify) for testing
- [captain](github.com/cyberhck/captain) for monitoring cron jobs
- [jwt-go](github.com/dgrijalva/jwt-go) for handling Auth using JWT
- [go-slack](github.com/multiplay/go-slack) for logging to slack
- [mock](github.com/golang/mock) for mocking interfaces
- [pushy](github.com/cyberhck/pushy) for sending push notifications to devices

### Ready to use?
As of v0.0.1, this boilerplate is ready for anyone to adapt,
modify and use. But beaware, that it's upto you to make this production
ready and usable for your needs.

### Upcoming changes
There WILL be breaking changes for a while till this gets to stable,
production ready state, till then, you can watch changes closely and adapt to it.

### Getting Started
As of now, it's not an application, but a package, it's because app engine requires you to do so.
This doesn't depend on app engine, but because of it's popularity, it made sense to declare as a package
which can be invoked really easily.

If you don't plan on going that route, you can always modify package name to your needs.

If you'd just like to see this working, follow the following steps:
- make sure you have docker and docker-compose
- make sure you've got go toolchain installed (older versions work, but they tend not to ignore `vendor` directory, so better to use newer go versions
- Then follow the following commands
```bash
cd $GOPATH/
mkdir -p src/github.com/fossapps/starter
cd src/github.com/fossapps/starter
git clone git@github.com:fossapps/golang_starter.git . # or you can use http url if you'd like
```
- then `cd` into starter directory
- run `cp .env.example .env`
- run `dep ensure`
- run `docker-compose up`
- run `make serve`

At this point, you'll have a API server running on port 8080 feel free to experiment.

### More Documentation
If you decide to use this for your needs, there are a few things which
you might have to keep in mind while using this starter project.

Visit WiKi to know how to move forward.

### Contributing
If and when you encounter some pain points or bugs, please feel free to open a new issue

See contributing.md to get started with contributions.

### License (MIT)
Visit license to see if this fits your needs.
