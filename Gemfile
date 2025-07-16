source 'https://rubygems.org'
git_source(:en_route) { |name| "https://bitbucket.org/enroute-mobi/#{name}.git" }

group :test do
  gem 'cucumber'
  gem 'rspec-expectations'
  gem 'mime-types'
  gem 'netrc'
  gem 'http-cookie'
  gem 'rest-client'
  gem 'pg'
  gem 'gtfs-rt', en_route: 'gtfs-rt'
  gem 'pry'
  gem 'siri-xsd', en_route: 'siri-xsd'
  gem 'activesupport'
  gem 'webrick'
  gem 'ara', en_route: 'ara-ruby', branch: 'ARA-1730-create-facility-model'
end

group :development do
  gem 'rake'
  gem 'license_finder'
  gem 'bundler-audit'
end
