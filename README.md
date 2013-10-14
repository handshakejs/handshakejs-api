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

### Heroku

```bash
git clone https://github.com/scottmotte/emailauth.git
cd emailauth
heroku create emailauth
git push heroku master
heroku run rake db:migrate
```

Next, create your first app.

```bash
curl -X POST https://emailauth.herokuapp.com/api/v0/apps.json

## Summary

### API Endpoint

* https://emailauth.herokuapp.com/api/v0

## Apps

To start using the emailauth.io API, you must first create an app.

### POST /apps/create

Pass an email and app_name to create your app at emailauth.herokuapp.com.

#### Definition

```bash
POST https://emailauth.herokuapp.com/api/v0/apps/create.json
```

#### Required Parameters

* email
* app_name

#### Example Request

```bash
curl -X POST https://emailauth.herokuapp.com/api/v0/apps/create.json -d "email=test@example.com" -d "app_name=myapp" -d "app_name=your_app_name"
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
POST https://emailauth.herokuapp.com/api/v0/login/request.json
```

#### Required Parameters

* email
* app_name

#### Example Request

```bash
curl -X POST https://mailauth.herokuapp.com/api/v0/login/request.json \ 
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
POST https://emailauth.herokuapp.com/api/v0/login/confirm.json
```

#### Required Parameters

* email
* authcode
* app_name

#### Example Request

```bash
curl -X POST https://emailauth.herokuapp.com/api/v0/login/confirm.json \
-d "email=test@example.com" \
-d "authcode=7389" \ 
-d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true,
  identity: {
    id: "idnt_1234348347834",  
    email: "test@example.com",
  }
}
```

## Identities

Under construction

## JS

The JS needs to build a form, submit via that form and send to /login/request. 

Then it needs to show the authcode portion of the form.

Then when the person puts in their auth code, it needs to send to /login/confirm to get a success response back.

If a success, then the developer programs her bit of code to log the user in via a session. How will she do this? Usually a session is generated on the back-end? Wouldn't it then make it easy for anyone to create a session out of thin air?

What's usually in a session - a hashed version of their email address that you do a lookup on? I need to look into how rails session end up working. It seems like the site then does a lookup on the user id via the session that was created. Usually that session is the email key plus a salt and then hashed. Then it does a reverse lookup on that internally in the app every time someone bounces around the web pages. 

Based on that they would just make a web request to generate the session via ajax - using the email that was passed as valid. The problem is - is then they are just passing an email and anyone could fake that with anyone's email. The token needs to be generated on the back-end probably via a webhook? Gosh, webhooks would suck here becuase you might get stuck waiting. Instead internally they need to verify that the auth was right after submitting the authcode form. This would go internally to their system,a nd their system would make the call out to /login/confirmed. If successful then they could generate the magic session from there.

Or they could use a token setup. The authcode would work, and they'd get a success back and get the token out. Essentially then they are tasked with at least internally writing some bit of code to authenticate against. See Paypal auth, facebook auth, github auth and more for how to do this portion of it.

I could offer both webhook and program their own?? How would webhook work. It'd hit a url that say that the person is authenticated. This would do a lookup on the database via a key on the people table. So this webhook would literally make the call to generate that session key. The webhook would be quite sloppy.

I need to show them steps to setup their own route called /login/confirm and inside that it would make the request out to my software, see if it is correct, and then they'd do their own bit of code to log the user in with a session.

1. Show form
2. on submit send ajax call for /login/request
3. show authcode portion of form to user
4. user submits authcode portion and it goes to an interal /login/confirm url
5. Inside that internal route is a remote call to the idenitty.io/login/confirm.json url
6. If that comes back as a success then internally set your session and do all your further lookups after the fact with that. This is up to the builder of the software how they want to handle session lookup for their users.
