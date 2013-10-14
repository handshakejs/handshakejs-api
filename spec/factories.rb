FactoryGirl.define do
  factory :app do |f|
    f.email       "person1@mailinator.com"
    f.app_name    "myapp"
  end

  factory :login do |f|
    f.email "login1@mailinator.com"
    f.association(:app, factory: :app)
  end

  factory :identity do |f|
    f.email "login1@mailinator.com"
    f.association(:app, factory: :app)
  end
end
