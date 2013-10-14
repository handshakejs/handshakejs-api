class App < ActiveRecord::Base
  has_many :logins

  validates :email, presence: true, uniqueness: true 
  validates :app_name, presence: true, uniqueness: true

  before_save       :set_api_keys

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
end
