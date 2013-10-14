class App < ActiveRecord::Base
  has_many :logins
  has_many :identities

  validates :email, presence: true
  validates :app_name, presence: true, uniqueness: true, format: { with: /\A[a-z0-9]+\z/ }

  before_save       :set_api_keys
  before_save       :set_short_id

  private

  def set_api_keys
    self.secret_api_key         = random_secret_api_key if !secret_api_key
    self.public_api_key         = random_public_api_key if !public_api_key
  end

  def random_secret_api_key
    ["sk_", SecureRandom.hex(15)].join
  end

  def random_public_api_key
    ["pk_", SecureRandom.hex(15)].join
  end

  def set_short_id
    self.short_id = short_id || ["APP_", Shorty.next_short_id(App)].join
  end
end
