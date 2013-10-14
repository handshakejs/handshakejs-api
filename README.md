# [emailauth.io](https://emailauth.herokuapp.com) API Documentation

**API platform for authenticating users without requiring a password.**

The [emailauth.io](https://emailauth.herokuapp.com) API is based around REST. It uses standard HTTP authentication. [JSON](https://www.json.org/) is returned in all responses from the API, including errors.

I've tried to make it as easy to use as possible, but if you have any feedback please [let me know](mailto:scott@scottmotte.com).

* [Summary](#summary)
* [Apps](#apps)
* [Logins](#logins)
* [Identity](#identities)  personas, aliases, pseudonym, nom de plume, noms, pen names, plumes, 

alias.io, nom de plume, identies.io, identmail.io

## Installation

### Local

```bash
bundle
touch .env.development
```

In .env.development put your local postgres database url.

```bash
DATABASE_URL="postgres://scottmotte@localhost/emailauth_development"
```

Then run the app.

```bash
bundle exec rake db:migrate
bundle exec foreman start
```

## Summary

### API Endpoint

* https://emailauth.herokuapp.com/api

## Apps

To start using the emailauth.io API, you must first create an app.

### POST /apps

Pass an email and app_name to create your app at emailauth.herokuapp.com.

#### Definition

```bash
POST https://emailauth.herokuapp.com/api/apps.json
```

#### Required Parameters

* email
* app_name

#### Example Request

```bash
curl -X POST https://emailauth.herokuapp.com/api/apps.json -d "email=test@example.com" -d "app_name=myapp" -d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true,
  app: {
    id: "1",
    email: "test@example.com",
    app_name: "myapp"
  }
}
```

## Logins

### POST /login/request

Request a login.

#### Definition

```bash
POST https://emailauth.herokuapp.com/api/login/request.json
```

#### Required Parameters

* email
* app_name

#### Example Request

```bash
curl -X POST https://mailauth.herokuapp.com/api/login/request.json \ 
-d "email=test@example.com" \
-d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true
}
```

### POST /login/confirm

Confirm a login. Email and authcode must match to get a success response back. 

#### Definition

```bash
POST https://emailauth.herokuapp.com/api/login/confirm.json
```

#### Required Parameters

* email
* authcode
* app_name

#### Example Request

```bash
curl -X POST https://emailauth.herokuapp.com/api/login/confirm.json \
-d "email=test@example.com" \
-d "authcode=7389" \ 
-d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true
}
```

## Identities
