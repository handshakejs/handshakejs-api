# [Handshake](https://handshakejs.herokuapp.com) API Documentation

![](https://rawgithub.com/scottmotte/handshake-js/master/handshakejs.svg)

**API platform for authenticating users without requiring a password.**

## Installation

### Heroku

```bash
git clone https://github.com/scottmotte/handshake.git
cd handshake
heroku create handshakejs
heroku addons:add sendgrid
heroku addons:add redistogo
git push heroku master
heroku config:set FROM=login@yourapp.com
```

Next, create your first app.

```bash
curl -X POST https://handshakejs.herokuapp.com/api/v0/apps/create.json \
-d "email=you@email.com" \
-d "app_name=your_app_name"
```

Nice, that's all it takes to get your authentication system running. Now let's plug that into our app using the embeddable JavaScript.

Place a script tag wherever you want the login form displayed.  

```html
<script src='/path/to/handshake.js' 
        data-app_name="your_app_name" 
        data-root_url="https://handshakejs.herokuapp.com"></script>
```

Get the latest [handshake.js here](https://github.com/scottmotte/handshake-js/blob/master/build/handshake.js). Replace the `data-app_name` with your own.

Next, bind to the handshake:login_confirm event to get the successful login data. This is where you would make an internal request to your application to set the session for the user.

```html
<script>
  handshake.script.addEventListener('handshake:login_confirm', function(e) {
    console.log(e.data);
    $.post("/login/success", {email: e.data.identity.email}, function(data) {
      window.location.href = "/dashboard";
    });    
  }, false); 
</script>
```

Then you'd setup a route in your app at /login/success to do something like this (setting the session). Here's an example in ruby and there is also a [full example ruby app](https://github.com/scottmotte/handshake-example-ruby).

```ruby
  post "/login/success" do
    session[:user] = params[:email]
    redirect "/dashboard"
  end
```

## API Overview

The [handshakejs.herokuapp.com](https://handshakejs.herokuapp.com) API is based around REST. It uses standard HTTP authentication. [JSON](https://www.json.org/) is returned in all responses from the API, including errors.

I've tried to make it as easy to use as possible, but if you have any feedback please [let me know](mailto:scott@scottmotte.com).

* [Summary](#summary)
* [Apps](#apps)
* [Login](#login)

## Summary

### API Endpoint

* https://handshakejs.herokuapp.com/api/v0

## Apps

To start using the handshake API, you must first create an app.

### POST /apps/create

Pass an email and app_name to create your app at handshakejs.herokuapp.com.

#### Definition

```bash
POST https://handshakejs.herokuapp.com/api/v0/apps/create.json
```

#### Required Parameters

* email
* app_name

#### Example Request

```bash
curl -X POST https://handshakejs.herokuapp.com/api/v0/apps/create.json \
-d "email=test@example.com" \
-d "app_name=myapp"
```

#### Example Response
```javascript
{
  success: true,
  app: {
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
POST https://handshakejs.herokuapp.com/api/v0/login/request.json
```

#### Required Parameters

* email
* app_name

#### Example Request

```bash
curl -X POST https://handshakejs.herokuapp.com/api/v0/login/request.json \ 
-d "email=test@example.com" \
-d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true,
  identity: {
    email: "test@example.com",
    app_name: "your_app_name",
    authcode_expired_at: "1382833591309"
  }
}
```

### POST /login/confirm

Confirm a login. Email and authcode must match to get a success response back. 

#### Definition

```bash
POST https://handshakejs.herokuapp.com/api/v0/login/confirm.json
```

#### Required Parameters

* email
* authcode
* app_name

#### Example Request

```bash
curl -X POST https://handshakejs.herokuapp.com/api/v0/login/confirm.json \
-d "email=test@example.com" \
-d "authcode=7389" \ 
-d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true,
  identity: {
    email: "test@example.com",
    app_name: "your_app_name",
    authcode: "7389"
  }
}
```

## Database Schema with Redis

apps - collection of keys with all the app_names in there. SADD

apps/myappname - hash with all the data in there. HSET or HMSET

apps/theappname/identities - collection of keys with all the identities' emails in there. SADD

apps/theappname/identities/emailaddress HSET or HMSET

