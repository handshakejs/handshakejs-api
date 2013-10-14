ENV["RACK_ENV"] ||= "development"

require 'bundler'
Bundler.setup
Bundler.require(:default, ENV["RACK_ENV"].to_sym)

# 
# load local env files if they exist
#
Dotenv.load ".env.#{Sinatra::Base.environment.to_s}", '.env'
require './config/post_body_content_type_parser'
require './config/constants'

require './app/application'

Dir["./lib/*.rb"].each { |f| require f }
Dir["./app/*.rb"].each { |f| require f }
Dir["./app/models/*.rb"].each { |f| require f }
