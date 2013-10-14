class CreateLogins < ActiveRecord::Migration
  def up
    create_table :logins do |t|
      t.integer   :app_id
      t.integer   :identity_id
      t.string    :email 
      t.string    :authcode 
      t.string    :status
      t.datetime  :created_at
      t.datetime  :updated_at
    end
  end

  def down
    drop_table :logins
  end
end
