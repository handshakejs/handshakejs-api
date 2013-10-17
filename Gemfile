source "https://rubygems.org"
ruby '1.9.3'

gem 'activesupport',                    require: 'active_support/core_ext/hash/keys'
gem 'dotenv'
gem 'foreman'
gem 'json'
gem 'jsonify',                          require: %w(jsonify jsonify/tilt)
gem 'pg'
gem 'pony'
gem 'pry'
gem 'racksh'
gem 'rake'
gem 'rack-cors',                        require: 'rack/cors'
gem 'sinatra',                          require: 'sinatra/base'
gem 'sinatra-contrib',                  require: %w(sinatra/multi_route sinatra/reloader)
gem 'state_machine'
gem 'sinatra-activerecord', '~> 1.2.2', require: 'sinatra/activerecord'
gem 'unicorn'

group :development, :test do
  gem 'therubyracer'
  gem 'rb-fsevent'
end

group :test do
  gem 'database_cleaner'
  gem 'factory_girl'
  gem 'rack-test',                      require: 'rack/test'
  gem 'rspec'
  gem 'webmock', '1.6.2',               require: false
end
