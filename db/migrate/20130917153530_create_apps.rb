class CreateApps < ActiveRecord::Migration
  def up
    create_table :apps do |t|
      t.string    :short_id
      t.string    :email 
      t.string    :app_name 
      t.string    :secret_api_key 
      t.string    :public_api_key 
      t.datetime  :created_at
      t.datetime  :updated_at
    end
  end

  def down
    drop_table :apps
  end
end
