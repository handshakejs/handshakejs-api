require File.dirname(__FILE__) + '/config/boot.rb'

use Rack::Static, :urls => ['/assets', '/images', '/javascripts', '/stylesheets', '/favicon.ico', '/apple-touch-icon.png', '/apple-touch-icon-precomposed.png'], :root => 'public'

run Rack::URLMap.new({
  "/"                 => Application,
  "/api/v0"           => ApiV0,
})
