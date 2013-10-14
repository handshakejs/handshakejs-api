FactoryGirl.define do
  factory :app do |f|
    f.email       "person1@mailinator.com"
    f.app_name    "myapp"
  end

  factory :login do |f|
    f.email "login1@mailinator.com"
    f.association(:app, factory: :app)
  end






  factory :person do |f|
    f.sequence(:email)      {|n| "person-#{n}@mailinator.com"}
    f.password              "password"
  end

  factory :document do |f|
    f.url                   "http://scottmotte.com/assets/resume.pdf"
    f.association(:person, factory: :person)
  end

  factory :page do |f|
    f.sort 1
    f.association(:document, factory: :document)
  end

  factory :event do |f|
    f.type_string             "document.created"
    f.object_attributes       { FactoryGirl.create(:person).event_attributes_to_hash }
  end  
  
  factory :webhook do |f|
    f.url                     "http://google.com"
    f.association(:person, factory: :person)
  end
end
