class CreateIdentities < ActiveRecord::Migration
  def up
    create_table :identities do |t|
      t.string    :short_id
      t.integer   :app_id
      t.string    :email 
      t.datetime  :created_at
      t.datetime  :updated_at
    end
  end

  def down
    drop_table :identities
  end
end
